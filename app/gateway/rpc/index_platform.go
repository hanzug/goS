package rpc

import (
	"context"
	"go.uber.org/zap"

	pb "github.com/hanzug/goS/idl/pb/index_platform"
)

// BuildIndex 建立索引的RPC调用
func BuildIndex(ctx context.Context, req *pb.BuildIndexReq) (resp *pb.BuildIndexResp, err error) {
	resp, err = IndexPlatformClient.BuildIndexService(ctx, req)
	if err != nil {
		zap.S().Error("BuildIndex-BuildIndexService", err)
		return
	}

	return
}
