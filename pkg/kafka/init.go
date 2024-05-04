package kafka

import (
	"github.com/IBM/sarama"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"

	"github.com/hanzug/goS/config"
)

var GobalKafka sarama.Client

func InitKafka() {
	zap.S().Info(logs.RunFuncName())
	con := sarama.NewConfig()
	con.Producer.Return.Successes = true
	kafkaClient, err := sarama.NewClient(config.Conf.Kafka.Address, con)
	zap.S().Info("connect kafka ok", err)
	if err != nil {
		return
	}
	GobalKafka = kafkaClient
}
