package kfk_register

import (
	"context"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"

	"github.com/hanzug/goS/consts"
	"github.com/hanzug/goS/pkg/kfk/consume"
)

// RunInvertedIndex 启动处理前向索引数据的进程。
func RunInvertedIndex(ctx context.Context) {
	zap.S().Info(logs.RunFuncName())

	// 调用 ForwardIndexKafkaConsume 函数从 Kafka 中消费前向索引数据。
	// 指定了 Kafka 主题、消费者组ID以及分区分配策略为轮询。
	err := consume.ForwardIndexKafkaConsume(ctx, consts.KafkaCSVLoaderTopic, consts.KafkaCSVLoaderGroupId, consts.KafkaAssignorRoundRobin)
	if err != nil {
		zap.S().Error("RunInvertedIndex-ForwardIndexKafkaConsume err :", err)
	}
}
