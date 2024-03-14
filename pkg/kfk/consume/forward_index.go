package consume

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"

	"github.com/hanzug/goS/app/index_platform/repository/db/dao"
	"github.com/hanzug/goS/config"
	"github.com/hanzug/goS/consts"
	"github.com/hanzug/goS/pkg/kfk"
	logs "github.com/hanzug/goS/pkg/logger"
	"github.com/hanzug/goS/repository/mysql/model"
	"github.com/hanzug/goS/types"
)

// ForwardIndexKafkaConsume 正排索引的消费建立
func ForwardIndexKafkaConsume(ctx context.Context, topic, group, assignor string) (err error) {
	zap.S().Info(logs.RunFuncName())

	keepRunning := true
	zap.S().Infof("Starting a new Sarama consumer")
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	// 设置一个消费组
	consumer := ForwardIndexConsumer{
		Ready: make(chan bool),
	}
	configK := kfk.GetDefaultConsumeConfig(assignor)
	cancelCtx, cancel := context.WithCancel(ctx)
	client, err := sarama.NewConsumerGroup(config.Conf.Kafka.Address, group, configK)
	if err != nil {
		zap.S().Errorf("Error creating consumer group woker: %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	zap.S().Infof("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-cancelCtx.Done():
			zap.S().Infof("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			zap.S().Infof("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		zap.S().Errorf("Error closing woker: %v", err)
	}

	return
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		zap.S().Infof("Resuming consumption")
	} else {
		client.PauseAll()
		zap.S().Infof("Pausing consumption")
	}

	*isPaused = !*isPaused
}

// Consumer Sarama消费者群体的消费者
type ForwardIndexConsumer struct {
	Ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *ForwardIndexConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *ForwardIndexConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 方法：这个方法是消息处理的主要逻辑。它从 claim.Messages() 通道中读取消息，然后根据消息的内容进行处理。在这个例子中，它将消息解析为一个 Document 对象，然后将这个对象保存到数据库中。
// 必须启动 ConsumerGroupClaim 的 Messages() 消费者循环。
// 一旦 Messages() 通道关闭，处理程序必须完成其处理循环并退出。
func (consumer *ForwardIndexConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	zap.S().Info(logs.RunFuncName())
	ctx := context.Background()
	task := &types.Task{
		Columns:    []string{"doc_id", "title", "body", "url"},
		BiTable:    "data",
		SourceType: consts.DataSourceCSV,
	}
	iDao := dao.NewInputDataDao(ctx)
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				zap.S().Infof("message channel was closed")
				return nil
			}

			if task.SourceType == consts.DataSourceCSV {
				doc := new(types.Document)
				_ = doc.UnmarshalJSON(message.Value)
				// TODO: 后续再开发starrocks
				// up.Push(&types.Data2Starrocks{
				// 	DocId: doc.DocId,
				// 	Url:   "",
				// 	Title: doc.Title,
				// 	Desc:  doc.Body,
				// 	Score: 0.00, // 评分
				// })
				inputData := &model.InputData{
					DocId:  doc.DocId,
					Title:  doc.Title,
					Body:   doc.Body,
					Url:    "",
					Score:  0.0,
					Source: task.SourceType,
				}
				err := iDao.CreateInputData(inputData)
				if err != nil {
					zap.S().Errorf("iDao.CreateInputData:%+v", err)
				}
			}

			//zap.S().Debugf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
