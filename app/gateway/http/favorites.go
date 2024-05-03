package http

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hanzug/goS/app/gateway/rpc"
	pb "github.com/hanzug/goS/idl/pb/favorite"
	"github.com/hanzug/goS/pkg/ctl"
)

// ListFavorite 处理获取收藏夹列表的请求
func ListFavorite(ctx *gin.Context) {
	var req pb.FavoriteListReq             // 定义请求结构体
	if err := ctx.Bind(&req); err != nil { // 绑定请求参数到结构体
		zap.S().Errorf("Bind:%v", err)                             // 日志记录错误
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误")) // 返回错误信息
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context()) // 获取用户信息
	if err != nil {
		zap.S().Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id                  // 设置用户ID
	r, err := rpc.FavoriteList(ctx, &req) // 调用RPC方法获取数据
	if err != nil {
		zap.S().Errorf("FavoriteList:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "FavoriteList RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r)) // 返回成功响应
}

// CreateFavorite 处理创建收藏夹的请求
func CreateFavorite(ctx *gin.Context) {
	var req pb.FavoriteCreateReq
	if err := ctx.ShouldBind(&req); err != nil {
		zap.S().Errorf("ShouldBind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		zap.S().Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteCreate(ctx, &req)
	if err != nil {
		zap.S().Errorf("FavoriteCreate:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "FavoriteCreateReq RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

// UpdateFavorite 处理更新收藏夹的请求
func UpdateFavorite(ctx *gin.Context) {
	var req pb.FavoriteCreateReq
	if err := ctx.Bind(&req); err != nil {
		zap.S().Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		zap.S().Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteCreate(ctx, &req)
	if err != nil {
		zap.S().Errorf("FavoriteCreate:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "UpdateFavorite RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

// DeleteFavorite 处理删除收藏夹的请求
func DeleteFavorite(ctx *gin.Context) {
	var req pb.FavoriteDeleteReq
	if err := ctx.Bind(&req); err != nil {
		zap.S().Errorf("req:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		zap.S().Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteDelete(ctx, &req)
	if err != nil {
		zap.S().Errorf("FavoriteDelete:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "DeleteFavorite RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

// ListFavoriteDetail 处理获取收藏夹详细列表的请求
func ListFavoriteDetail(ctx *gin.Context) {
	var req pb.FavoriteDetailListReq
	if err := ctx.Bind(&req); err != nil {
		zap.S().Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		zap.S().Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteDetailList(ctx, &req)
	if err != nil {
		zap.S().Errorf("FavoriteDetailList:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "FavoriteDetailList RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

// CreateFavoriteDetail 处理创建收藏夹详细信息的请求
func CreateFavoriteDetail(ctx *gin.Context) {
	var req pb.FavoriteDetailCreateReq
	if err := ctx.Bind(&req); err != nil {
		zap.S().Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		zap.S().Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteDetailCreate(ctx, &req)
	if err != nil {
		zap.S().Errorf("FavoriteDetailCreate:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "FavoriteDetailCreate RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

// DeleteFavoriteDetail 处理删除收藏夹详细信息的请求
func DeleteFavoriteDetail(ctx *gin.Context) {
	var req pb.FavoriteDetailDeleteReq
	if err := ctx.Bind(&req); err != nil {
		zap.S().Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	user, err := ctl.GetUserInfo(ctx.Request.Context())
	if err != nil {
		zap.S().Errorf("GetUserInfo:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "获取用户信息错误"))
		return
	}
	req.UserId = user.Id
	r, err := rpc.FavoriteDetailDelete(ctx, &req)
	if err != nil {
		zap.S().Errorf("FavoriteDetailDelete:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "FavoriteDetailDelete RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}
