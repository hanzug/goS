package kfk_register

import (
	"context"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"
)

func RegisterJob(ctx context.Context) {
	zap.S().Info(logs.RunFuncName())
	newCtx := ctx
	// go RunTireTree(newCtx) // TODO:这个有点问题，后续优化再开启
	// newCtx = ctx
	go RunInvertedIndex(newCtx)
}
