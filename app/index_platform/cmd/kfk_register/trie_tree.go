package kfk_register

import (
	"context"
	"fmt"
	"github.com/hanzug/goS/app/index_platform/kafka/consumer"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"

	"github.com/hanzug/goS/consts"
)

func RunTireTree(ctx context.Context) {
	zap.S().Info(logs.RunFuncName())
	// Trie处理
	err := consumer.TrieTreeKafkaConsume(ctx, consts.KafkaTrieTreeTopic, consts.KafkaTrieTreeGroupId, consts.KafkaAssignorRoundRobin)
	if err != nil {
		fmt.Println("RunTireTree-TrieTreeKafkaConsume :", err)
	}
}
