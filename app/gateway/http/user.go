package http

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hanzug/goS/app/gateway/rpc"
	pb "github.com/hanzug/goS/idl/pb/user"
	"github.com/hanzug/goS/pkg/ctl"
	"github.com/hanzug/goS/pkg/jwt"
	"github.com/hanzug/goS/types"
)

// UserRegister 用户注册
func UserRegister(ctx *gin.Context) {
	var userReq pb.UserRegisterReq
	if err := ctx.ShouldBind(&userReq); err != nil {
		zap.S().Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}
	r, err := rpc.UserRegister(ctx, &userReq)
	if err != nil {
		zap.S().Errorf("UserRegister:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "UserRegister RPC服务调用错误"))
		return
	}

	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, r))
}

// UserLogin 用户登录
func UserLogin(ctx *gin.Context) {
	var req pb.UserLoginReq
	if err := ctx.ShouldBind(&req); err != nil {
		zap.S().Errorf("Bind:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "绑定参数错误"))
		return
	}

	userResp, err := rpc.UserLogin(ctx, &req)
	if err != nil {
		zap.S().Errorf("RPC UserLogin:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "UserLogin RPC服务调用错误"))
		return
	}

	aToken, rToken, err := jwt.GenerateToken(userResp.UserDetail.UserId, userResp.UserDetail.UserName)
	if err != nil {
		zap.S().Errorf("RPC GenerateToken:%v", err)
		ctx.JSON(http.StatusOK, ctl.RespError(ctx, err, "加密错误"))
		return
	}
	uResp := &types.UserTokenData{
		User:         userResp.UserDetail,
		AccessToken:  aToken,
		RefreshToken: rToken,
	}
	ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, uResp))
}
