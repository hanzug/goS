package http

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hanzug/goS/app/gateway/rpc"
	pb "github.com/hanzug/goS/idl/pb/search_engine"
	"github.com/hanzug/goS/pkg/ctl"
)

// SearchEngineSearch 搜索
func SearchEngineSearch(ctx *gin.Context) {
	var req *pb.SearchEngineRequest
	if err := ctx.ShouldBind(&req); err != nil {
		zap.S().Errorf("SearchEngineSearch-ShouldBind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}

	r, err := rpc.SearchEngineSearch(ctx, req)
	if err != nil {
		zap.S().Errorf("SearchEngineSearch:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "SearchEngineSearch RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

// WordAssociation 词条联想
func WordAssociation(ctx *gin.Context) {
	var req *pb.SearchEngineRequest
	if err := ctx.ShouldBind(&req); err != nil {
		zap.S().Errorf("WordAssociation-ShouldBind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}

	r, err := rpc.WordAssociation(ctx, req)
	if err != nil {
		zap.S().Errorf("WordAssociation:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "WordAssociation RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}
