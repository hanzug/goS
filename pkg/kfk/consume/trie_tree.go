package consume

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"

	"github.com/hanzug/goS/app/index_platform/trie"
	"github.com/hanzug/goS/config"
	"github.com/hanzug/goS/pkg/kfk"
)

// TrieTreeKafkaConsume token词的消费建立
func TrieTreeKafkaConsume(ctx context.Context, topic, group, assignor string) (err error) {
	zap.S().Infof("Starting a new Sarama consumer")
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	// 设置一个消费组
	consumer := TrieTreeConsumer{
		Ready: make(chan bool),
	}

	configK := kfk.GetDefaultConsumeConfig(assignor)
	cancelCtx, cancel := context.WithCancel(ctx)
	client, err := sarama.NewConsumerGroup(config.Conf.Kafka.Address, group, configK)
	if err != nil {
		zap.S().Errorf("Error creating consumer group woker: %v", err)
	}

	go func() {
		for {
			if err = client.Consume(cancelCtx, []string{topic}, &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				zap.S().Errorf("Error from consumer: %v", err)
			}
			if cancelCtx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready
	cancel()

	return
}

// TrieTreeConsumer Sarama消费者群体的消费者
type TrieTreeConsumer struct {
	Ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *TrieTreeConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *TrieTreeConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 必须启动 ConsumerGroupClaim 的 Messages() 消费者循环。
// 一旦 Messages() 通道关闭，处理程序必须完成其处理循环并退出。
func (consumer *TrieTreeConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// ctx := context.Background()
	gapTime := 2 * time.Minute
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				zap.S().Infof("message channel was closed")
				return nil
			}
			// 构建trie tree树
			trie.GobalTrieTree.Insert(string(message.Value))
			// zap.S().Infof("TrieTreeConsumer Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		// https://github.com/IBM/sarama/issues/1192

		case <-time.After(gapTime):
			zap.S().Infof("ConsumeClaim starting store dict")
			// _ = storage.GlobalTrieDBs.StorageDict(trie.GobalTrieTree) // TODO:后续看看能不能实现一个全局的triedb，每次都先读取存量进行初始化，再插入增量...
			zap.S().Infof("ConsumeClaim ending store dict")

		case <-session.Context().Done():
			zap.S().Infof("TrieTreeConsumer Done!")
			return nil
		}
	}
}

// func mergeTrieTree(node string) {
// 	trie.GobalTrieTree.Insert(node)
// 	gapTime := 2 * time.Minute
// 	for {
// 		select {
// 		case <-time.After(gapTime):
// 			_ = storage.GlobalTrieDBs.StorageDict(trie.GobalTrieTree)
// 		}
// 	}
// }
