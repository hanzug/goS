package kfk_register

import (
	"context"
	"fmt"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"

	"github.com/hanzug/goS/consts"
	"github.com/hanzug/goS/pkg/kfk/consume"
)

func RunTireTree(ctx context.Context) {
	zap.S().Info(logs.RunFuncName())
	// Trie处理
	err := consume.TrieTreeKafkaConsume(ctx, consts.KafkaTrieTreeTopic, consts.KafkaTrieTreeGroupId, consts.KafkaAssignorRoundRobin)
	if err != nil {
		fmt.Println("RunTireTree-TrieTreeKafkaConsume :", err)
	}
}
