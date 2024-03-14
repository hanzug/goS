package ranking

import (
	"sort"

	"github.com/hanzug/goS/pkg/util/relevant"
	"github.com/hanzug/goS/types"
)

// CalculateScoreBm25 计算给定令牌在一系列搜索项中的BM25得分，并根据得分对搜索项进行排序。
func CalculateScoreBm25(token string, searchItem []*types.SearchItem) (resp []*types.SearchItem) {
	// 初始化一个字符串切片来存储每个搜索项的内容。
	contents := make([]string, 0)

	// 遍历搜索项，将每个项的内容添加到contents切片中。
	for i := range searchItem {
		contents = append(contents, searchItem[i].Content)
	}

	// 根据contents切片中的内容创建一个语料库。
	corpus, _ := relevant.MakeCorpus(contents)

	// 使用语料库为每个内容项创建文档表示。
	docs := relevant.MakeDocuments(contents, corpus)
	// 初始化一个新的TF-IDF计算器。
	tf := relevant.New()
	// 将每个文档添加到TF-IDF计算器中。
	for _, doc := range docs {
		tf.Add(doc)
	}
	// 为语料库中的每个术语计算逆文档频率（IDF）。
	tf.CalculateIDF()
	// 为搜索令牌创建一个文档表示。
	tokenRecall := relevant.Doc{corpus[token]}
	// 计算相对于搜索令牌的每个文档的BM25得分。
	bm25Scores := relevant.BM25(tf, tokenRecall, docs, 1.5, 0.75)
	// 按降序对BM25得分进行排序。
	sort.Sort(sort.Reverse(bm25Scores))

	// 根据BM25得分更新每个搜索项的得分，并过滤掉得分为0的项。
	for i := range bm25Scores {
		if bm25Scores[i].Score == 0.0 {
			continue
		}
		searchItem[bm25Scores[i].ID].Score = bm25Scores[i].Score
	}
	// 根据得分对搜索项进行降序排序。
	sort.Slice(searchItem, func(i, j int) bool {
		return searchItem[i].Score > searchItem[j].Score
	})

	// 将排序后的搜索项赋值给响应变量。
	resp = searchItem

	return
}
