package http

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hanzug/goS/app/gateway/rpc"
	pb "github.com/hanzug/goS/idl/pb/index_platform"
	"github.com/hanzug/goS/pkg/ctl"
)

func BuildIndexByFiles(ctx *gin.Context) {
	var req pb.BuildIndexReq
	if err := ctx.ShouldBind(&req); err != nil {
		zap.S().Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}

	r, err := rpc.BuildIndex(ctx, &req)
	if err != nil {
		zap.S().Errorf("BuildIndexByFiles:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "BuildIndexByFiles RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}
