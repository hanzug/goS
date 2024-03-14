package rpc

import (
	"context"
	"go.uber.org/zap"

	"github.com/hanzug/goS/idl/pb/mapreduce"
)

// MasterAssignTask 通过 master 发送任务
func MasterAssignTask(ctx context.Context, taskReq *mapreduce.MapReduceTask) (resp *mapreduce.MapReduceTask, err error) {
	resp, err = MapReduceClient.MasterAssignTask(ctx, taskReq)
	if err != nil {
		zap.S().Error("MasterAssignTask-MapReduceClient", err)
		return
	}

	return
}

// MasterTaskCompleted 通知 master 任务完成的RPC调用
func MasterTaskCompleted(ctx context.Context, task *mapreduce.MapReduceTask) (resp *mapreduce.MasterTaskCompletedResp, err error) {
	resp, err = MapReduceClient.MasterTaskCompleted(ctx, task)
	if err != nil {
		zap.S().Error("MasterTaskCompleted-MapReduceClient", err)
		return
	}

	return
}
