package recall

import (
	"context"
	"go.uber.org/zap"

	"github.com/RoaringBitmap/roaring"
	"github.com/samber/lo"

	"github.com/hanzug/goS/app/search_engine/analyzer"
	"github.com/hanzug/goS/app/search_engine/ranking"
	"github.com/hanzug/goS/app/search_engine/repository/db/dao"
	"github.com/hanzug/goS/app/search_engine/repository/storage"
	"github.com/hanzug/goS/repository/redis"
	"github.com/hanzug/goS/types"
)

// Recall 查询召回
type Recall struct {
}

func NewRecall() *Recall {
	return &Recall{}
}

// Search 入口
func (r *Recall) Search(ctx context.Context, query string) (res []*types.SearchItem, err error) {
	splitQuery, err := analyzer.GseCutForRecall(query)
	if err != nil {
		zap.S().Errorf("text2postingslists err: %v", err)
		return
	}

	res, err = r.searchDoc(ctx, splitQuery)

	return
}

// SearchQuery 入口
func (r *Recall) SearchQuery(query string) (resp []string, err error) {
	dictTreeList := make([]string, 0, 1e3)
	for _, trieDb := range storage.GlobalTrieDB {
		// 获取 Trie 树：对于全局 Trie 数据库 GlobalTrieDB 中的每一个 Trie 数据库，这个函数首先获取这个数据库的 Trie 树。
		trie, errx := trieDb.GetTrieTreeDict()
		if errx != nil {
			zap.S().Error(errx)
			continue
		}
		//查找匹配的字符串：然后，这个函数使用 FindAllByPrefixForRecall 函数查找所有以 query 为前缀的字符串。这些字符串被添加到 dictTreeList 列表中。
		queryTrie := trie.FindAllByPrefixForRecall(query)
		dictTreeList = append(dictTreeList, queryTrie...)
	}
	//去重
	resp = lo.Uniq(dictTreeList)
	return
}

func (r *Recall) searchDoc(ctx context.Context, tokens []string) (recalls []*types.SearchItem, err error) {
	recalls = make([]*types.SearchItem, 0)
	allPostingsList := []*types.PostingsList{}
	for _, token := range tokens {
		//获取倒排列表
		docIds, errx := redis.GetInvertedIndexTokenDocIds(ctx, token)
		//全局倒排列表
		var postingsList []*types.PostingsList
		if errx != nil || docIds == nil {
			// 如果缓存不存在，就去索引表里面读取
			postingsList, err = fetchPostingsByToken(token)
			if err != nil {
				zap.S().Error(err)
				continue
			} else {
				// todo 缓存一致性优化（设置为可选）
				// 如果缓存存在，就直接读缓存，不用担心实时性问题，缓存10分钟清空一次，这延迟是能接受到
				postingsList = append(postingsList, &types.PostingsList{
					Term:   token,
					DocIds: docIds,
				})
			}
		}
		allPostingsList = append(allPostingsList, postingsList...)
	}

	// 排序打分
	iDao := dao.NewInputDataDao(ctx)
	for _, p := range allPostingsList {
		if p == nil || p.DocIds == nil || p.DocIds.IsEmpty() {
			continue
		}
		recallData, _ := iDao.ListInputDataByDocIds(p.DocIds.ToArray())
		searchItems := ranking.CalculateScoreBm25(p.Term, recallData)
		recalls = append(recalls, searchItems...)
	}

	zap.S().Infof("recalls size:%v", len(recalls))

	return
}

// 获取 token 所有seg的倒排表数据
func fetchPostingsByToken(token string) (postingsList []*types.PostingsList, err error) {
	// 遍历存储index的地方，token对应的doc Id 全部取出
	postingsList = make([]*types.PostingsList, 0, 1e6)
	for _, inverted := range storage.GlobalInvertedDB {
		docIds, errx := inverted.GetInverted([]byte(token))
		if errx != nil {
			zap.S().Info(errx)
			continue
		}
		output := roaring.New()
		_ = output.UnmarshalBinary(docIds)
		// 存放到数组当中
		postings := &types.PostingsList{
			Term:   token,
			DocIds: output,
		}
		postingsList = append(postingsList, postings)
	}

	return
}
