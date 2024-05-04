package consumer

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
	"github.com/hanzug/goS/pkg/kafka"
	logs "github.com/hanzug/goS/pkg/logger"
	"github.com/hanzug/goS/repository/mysql/model"
	"github.com/hanzug/goS/types"
)

// ForwardIndexKafkaConsume 正排索引的消费建立
func ForwardIndexKafkaConsume(ctx context.Context, topic, group, assignor string) (err error) {

	zap.S().Info(logs.RunFuncName())

	// 标志变量，用于控制消费循环
	keepRunning := true

	zap.S().Infof("Starting a new Sarama consumer")
	// 设置Sarama的日志输出到标准输出
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	// 初始化消费者组对象
	consumer := ForwardIndexConsumer{
		Ready: make(chan bool),
	}
	// 获取默认的消费者配置
	configK := kafka.GetDefaultConsumeConfig(assignor)
	// 创建一个可取消的上下文
	cancelCtx, cancel := context.WithCancel(ctx)
	// 尝试创建消费者组
	client, err := sarama.NewConsumerGroup(config.Conf.Kafka.Address, group, configK)
	if err != nil {
		// 创建失败，记录错误
		zap.S().Errorf("Error creating consumer group worker: %v", err)
	}

	// 消费暂停标志
	consumptionIsPaused := false

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// 消费消息
			if err = client.Consume(cancelCtx, []string{topic}, &consumer); err != nil {
				// 如果消费者组被关闭，退出循环
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				// 记录消费过程中的错误
				zap.S().Errorf("Error from consumer: %v", err)
			}

			zap.S().Info("消费正排索引")
			// 如果上下文被取消，退出循环
			if cancelCtx.Err() != nil {
				return
			}
			// 重置Ready通道，准备下一轮消费
			consumer.Ready = make(chan bool)
		}
	}()

	// 等待消费者准备就绪
	<-consumer.Ready
	zap.S().Infof("Sarama consumer up and running!...")
	// 设置信号通道，用于接收中断信号
	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, os.Interrupt)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// 循环，直到接收到终止信号
	for keepRunning {
		select {
		case <-cancelCtx.Done():
			// 上下文被取消，记录并退出循环
			zap.S().Infof("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			// 接收到终止信号，记录并退出循环
			zap.S().Infof("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			// 接收到用户定义信号，切换消费流状态
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	// 取消上下文
	cancel()
	// 等待所有协程完成
	wg.Wait()
	// 尝试关闭消费者组并处理错误
	if err = client.Close(); err != nil {
		zap.S().Errorf("Error closing worker: %v", err)
	}

	return
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	zap.S().Info(logs.RunFuncName())
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
	zap.S().Info(logs.RunFuncName())
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *ForwardIndexConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	zap.S().Info(logs.RunFuncName())
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
