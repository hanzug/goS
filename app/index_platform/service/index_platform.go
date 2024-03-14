package service

import (
	"context"
	"fmt"
	"github.com/hanzug/goS/pkg/clone"
	"go.uber.org/zap"
	"hash/fnv"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/RoaringBitmap/roaring"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/spf13/cast"

	"github.com/hanzug/goS/app/index_platform/analyzer"
	"github.com/hanzug/goS/app/index_platform/input_data"
	"github.com/hanzug/goS/app/index_platform/repository/storage"
	cconsts "github.com/hanzug/goS/consts"
	"github.com/hanzug/goS/consts/e"
	pb "github.com/hanzug/goS/idl/pb/index_platform"
	"github.com/hanzug/goS/pkg/mapreduce"
	"github.com/hanzug/goS/pkg/timeutils"
	"github.com/hanzug/goS/pkg/trie"
	"github.com/hanzug/goS/types"
)

type IndexPlatformSrv struct {
	*pb.UnimplementedIndexPlatformServiceServer
}

var (
	IndexPlatIns  *IndexPlatformSrv
	IndexPlatOnce sync.Once
)

func GetIndexPlatformSrv() *IndexPlatformSrv {
	IndexPlatOnce.Do(func() {
		IndexPlatIns = new(IndexPlatformSrv)
	})
	return IndexPlatIns
}

// BuildIndexService 构建索引
func (s *IndexPlatformSrv) BuildIndexService(ctx context.Context, req *pb.BuildIndexReq) (resp *pb.BuildIndexResp, err error) {
	// 时间估计
	resp = new(pb.BuildIndexResp)
	resp.Code = e.SUCCESS
	resp.Message = e.GetMsg(e.SUCCESS)
	invertedIndex := cmap.New[*roaring.Bitmap]() // 倒排索引
	dictTrie := trie.NewTrie()                   // 前缀树

	zap.S().Infof("BuildIndexService Start req: %v", req.FilePath)
	// mapreduce 这个是用chan和goroutine来代替master和worker的rpc调用，避免了频繁的rpc调用
	_, _ = mapreduce.MapReduce(func(source chan<- []byte) {
		// 读取http传来的文件地址
		zap.S().Info("1mapreduce 读取文件")
		for _, path := range req.FilePath {
			content, _ := os.ReadFile(path)
			source <- content
		}
	}, func(item []byte, writer mapreduce.Writer[[]*types.KeyValue], cancel func(error)) {
		// 控制并发
		zap.S().Info("2mapreduce map阶段启动")
		var wg sync.WaitGroup
		ch := make(chan struct{}, 3)

		keyValueList := make([]*types.KeyValue, 0, 1e3)
		lines := strings.Split(string(item), "\r\n")
		for _, line := range lines[1:] {
			ch <- struct{}{}

			zap.S().Info("3mapreduce map函数执行")

			wg.Add(1)
			docStruct, _ := input_data.Doc2Struct(line) // line 转 doc struct
			if docStruct.DocId == 0 {
				continue
			}

			// 分词
			tokens, _ := analyzer.GseCutForBuildIndex(docStruct.DocId, docStruct.Body)
			for _, v := range tokens {
				if v.Token == "" || v.Token == " " {
					continue
				}
				keyValueList = append(keyValueList, &types.KeyValue{Key: v.Token, Value: cast.ToString(v.DocId)})
				zap.S().Info("4插入TrieTree")
				dictTrie.Insert(v.Token)
			}

			// 建立正排索引
			go func(docStruct *types.Document) {
				zap.S().Info("5发送到kafka")
				err = input_data.DocData2Kfk(docStruct)
				if err != nil {
					zap.S().Error(err)
				}
				defer wg.Done()
				<-ch
			}(docStruct)
		}
		wg.Wait()

		// // 构建前缀树 // TODO:kafka处理
		// go func(tokenList []string) {
		// 	err = input_data.DocTrie2Kfk(tokenList)
		// 	if err != nil {
		// 		zap.S().Error("DocTrie2Kfk", err)
		// 	}
		// }(tokenList)

		// shuffle 排序过程
		sort.Sort(types.ByKey(keyValueList))
		writer.Write(keyValueList)
	}, func(pipe <-chan []*types.KeyValue, writer mapreduce.Writer[string], cancel func(error)) {
		for values := range pipe {
			for _, v := range values { // 构建倒排索引
				if value, ok := invertedIndex.Get(v.Key); ok {
					value.AddInt(cast.ToInt(v.Value))
					invertedIndex.Set(v.Key, value)
				} else {
					docIds := roaring.NewBitmap()
					docIds.AddInt(cast.ToInt(v.Value))
					invertedIndex.Set(v.Key, docIds)
				}
			}
		}
	})

	var wg sync.WaitGroup

	go func() {
		newCtx := clone.NewContextWithoutDeadline()
		newCtx.Clone(ctx)
		err = storeInvertedIndexByHash(newCtx, invertedIndex)
		if err != nil {
			zap.S().Error("storeInvertedIndexByHash error ", err)
		}
	}()

	zap.S().Info("storeInvertedIndexByHash End")
	go func() {
		newCtx := clone.NewContextWithoutDeadline()
		newCtx.Clone(ctx)
		err = storeDictTrieByHash(newCtx, dictTrie)
		if err != nil {
			zap.S().Error("storeDictTrieByHash error ", err)
		}
	}()
	wg.Wait()

	return
}

// storeInvertedIndexByHash 分片存储
func storeInvertedIndexByHash(ctx context.Context, invertedIndex cmap.ConcurrentMap[string, *roaring.Bitmap]) (err error) {
	dir, _ := os.Getwd()
	outName := fmt.Sprintf("%s/%s.%s", dir, timeutils.GetTodayDate(), cconsts.InvertedBucket)
	invertedDB := storage.NewInvertedDB(outName)
	// 找出所有的key进行存储
	for k, val := range invertedIndex.Items() {
		outByte, errx := val.MarshalBinary()
		if errx != nil {
			zap.S().Error("storeInvertedIndexByHash-MarshalBinary", errx)
			continue
		}
		err = invertedDB.StoragePostings(k, outByte)
		if err != nil {
			zap.S().Error("storeInvertedIndexByHash-StoragePostings", err)
			continue
		}
	}

	err = redis.PushInvertedPath(ctx, redis.InvertedIndexDbPathKey, []string{outName})
	zap.S().Info("______________redis PushInvertedPath")
	str, _ := redis.ListInvertedPath(ctx, redis.InvertedIndexDbPathKey)
	zap.S().Info("", zap.Any("_______________redis get", str))
	fmt.Println(outName)
	if err != nil {
		zap.S().Error(err)
		return
	}

	// TODO: hash 分片存储
	// dir, _ := os.Getwd()
	// keys := invertedIndex.Keys()
	// buffer := make([][]*types.KeyValue, consts.ReduceDefaultNum)
	// for i, v := range keys {
	// 	val, _ := invertedIndex.Get(v)
	// 	slot := iHash(v) % consts.ReduceDefaultNum
	// 	buffer[slot] = append(buffer[slot])
	// 	fmt.Println(v, val)
	// }
	// outName := fmt.Sprintf("%s/%d.%s", dir, i, cconsts.InvertedBucket)

	return
}

// storeInvertedIndexByHash 分片存储
func storeDictTrieByHash(ctx context.Context, dict *trie.Trie) (err error) {
	// TODO: 抽离一个hash存储的方法
	dir, _ := os.Getwd()
	outName := fmt.Sprintf("%s/%s.%s", dir, timeutils.GetTodayDate(), cconsts.TrieTreeBucket)
	trieDB := storage.NewTrieDB(outName)
	err = trieDB.StorageDict(dict)
	if err != nil {
		zap.S().Error(err)
		return
	}

	err = redis.PushInvertedPath(ctx, redis.TireTreeDbPathKey, []string{outName})
	if err != nil {
		zap.S().Error(err)
		return
	}

	return
}

// iHash 哈希作用
func iHash(key string) int64 { // nolint:golint,unused
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int64(h.Sum32() & 0x7fffffff)
}
