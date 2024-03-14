package relevant

import (
	"sort"

	"github.com/xtgo/set"

	"github.com/hanzug/goS/app/search_engine/analyzer"
)

// DocScore 结构体：这个结构体表示一个文档的得分，其中 ID 是文档的唯一标识符，Score 是该文档的得分。
type DocScore struct {
	ID    int
	Score float64
}

// DocScores 类型：这是一个 DocScore 的切片，它实现了 sort.Interface 接口，所以它可以被 sort 包中的函数排序。
type DocScores []DocScore

func (ds DocScores) Len() int           { return len(ds) }
func (ds DocScores) Less(i, j int) bool { return ds[i].Score < ds[j].Score }
func (ds DocScores) Swap(i, j int) {
	ds[i].Score, ds[j].Score = ds[j].Score, ds[i].Score
	ds[i].ID, ds[j].ID = ds[j].ID, ds[i].ID
}

// BM25 函数：这是一个实现了 BM25 算法的函数，BM25 是一种信息检索中常用的评分函数。
// 它接受一个查询文档和一组文档，然后返回每个文档的得分。这个函数首先将查询文档转换为词袋模型，然后对每个文档计算其得分。
// k1 should be between 1.2 and 2.
// b should be around 0.75
func BM25(tf *TFIDF, query Document, docs []Document, k1, b float64) DocScores {
	// 词袋模型
	q := BOW(query)
	// 创建一个和查询文档同样大小的工作向量
	w := make([]int, len(q))
	// 将查询文档复制到工作向量
	copy(w, q)
	// 计算文档的平均长度
	avgLen := float64(tf.Len) / float64(tf.Docs)

	// 创建一个用于存储每个文档得分的切片
	scores := make([]float64, 0, len(docs))
	for _, doc := range docs {
		// 将文档转换为词袋模型
		d := BOW(doc)
		// 将文档添加到工作向量
		w = append(w, d...)
		// 计算查询文档和当前文档的交集的大小
		size := set.Inter(sort.IntSlice(w), len(q))
		// 获取查询文档和当前文档的交集
		n := w[:size]

		// 创建一个用于存储当前文档得分的切片
		score := make([]float64, 0, len(n))
		// 获取当前文档的长度
		docLen := float64(len(d))
		for _, id := range n {
			// 计算得分的分子
			num := tf.TF[id] * (k1 + 1)
			// 计算得分的分母
			denom := tf.TF[id] + k1*(1-b+b*docLen/avgLen)
			// 获取当前词的逆文档频率
			idf := tf.IDF[id]
			// 计算得分并添加到得分切片
			score = append(score, idf*num/denom)
		}
		// 计算当前文档的总得分并添加到得分切片
		scores = append(scores, sum(score))

		// 重置工作向量
		copy(w, q)
		w = w[:len(q)]
	}
	// 创建一个用于返回的 DocScores
	var retVal DocScores
	for i := range docs {
		// 将每个文档的 ID 和得分添加到 retVal
		retVal = append(retVal, DocScore{i, scores[i]})
	}
	return retVal
}

func sum(a []float64) float64 {
	var retVal float64
	for _, f := range a {
		retVal += f
	}
	return retVal
}

type Doc []int

func (d Doc) IDs() []int { return d }

// MakeCorpus 从一系列文本中创建一个语料库
// 返回两个结果：一个映射，将每个词汇映射到一个唯一的整数标识符；一个切片，将每个整数标识符映射回其对应的词汇。
func MakeCorpus(a []string) (map[string]int, []string) {
	retVal := make(map[string]int) // 创建一个映射
	invRetVal := make([]string, 0) // 创建一个切片
	var id int
	for _, s := range a { // 遍历每个文本
		tokens, _ := analyzer.GseCutForRecall(s) // 对文本进行分词
		for _, f := range tokens {               // 遍历每个词汇
			if _, ok := retVal[f]; !ok { // 如果这个词汇还没有被映射到一个整数
				retVal[f] = id                   // 将这个词汇映射到一个新的整数
				invRetVal = append(invRetVal, f) // 将这个整数映射回这个词汇
				id++                             // 整数加一
			}
		}
	}
	return retVal, invRetVal // 返回映射和切片
}

// MakeDocuments 从一系列文本中创建一个文档集合，每个文档是一个整数切片，这些整数是词汇的唯一标识符。
func MakeDocuments(a []string, c map[string]int) []Document {
	retVal := make([]Document, 0, len(a))
	for _, s := range a {
		var ts []int
		tokens, _ := analyzer.GseCutForRecall(s)
		for _, f := range tokens {
			id := c[f]
			ts = append(ts, id)
		}
		retVal = append(retVal, Doc(ts))
	}
	return retVal
}
