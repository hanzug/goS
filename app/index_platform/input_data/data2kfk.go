package input_data

import (
	"github.com/hanzug/goS/consts"
	"github.com/hanzug/goS/pkg/kafka"
	logs "github.com/hanzug/goS/pkg/logger"
	"github.com/hanzug/goS/types"
	"go.uber.org/zap"
)

// DocData2Kfk Doc数据处理
func DocData2Kfk(doc *types.Document) (err error) {
	zap.S().Info(logs.RunFuncName())

	zap.S().Info("生产者生产一次正排索引")

	doctByte, _ := doc.MarshalJSON()
	err = kafka.KafkaProducer(consts.KafkaCSVLoaderTopic, doctByte)
	if err != nil {
		zap.S().Errorf("DocData2Kfk-KafkaCSVLoaderTopic :%+v", err)
		return
	}

	return
}

// DocTrie2Kfk Trie树构建
func DocTrie2Kfk(tokens []string) (err error) {
	zap.S().Info(logs.RunFuncName())

	zap.S().Info("生产者生产一次TrieTree")

	for _, k := range tokens {
		err = kafka.KafkaProducer(consts.KafkaTrieTreeTopic, []byte(k))
	}

	if err != nil {
		zap.S().Errorf("DocTrie2Kfk-KafkaTrieTreeTopic :%+v", err)
		return
	}

	return
}
