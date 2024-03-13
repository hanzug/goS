package service

import (
	"context"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"
	"sync"

	"github.com/hanzug/goS/app/user/internal/repository/db/dao"
	e2 "github.com/hanzug/goS/consts/e"
	pb "github.com/hanzug/goS/idl/pb/user"
)

// UserSrvIns grpc用户实例
var UserSrvIns *UserSrv

// UserSrvOnce 保证只执行一次
var UserSrvOnce sync.Once

type UserSrv struct {
	pb.UnimplementedUserServiceServer
}

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

func (u *UserSrv) UserLogin(ctx context.Context, req *pb.UserLoginReq) (resp *pb.UserDetailResponse, err error) {

	zap.S().Info(logs.RunFuncName())
	resp = new(pb.UserDetailResponse)
	resp.Code = e2.SUCCESS

	r, err := dao.NewUserDao(ctx).GetUserInfo(req)
	if err != nil {
		resp.Code = e2.ERROR
		return
	}
	resp.UserDetail = &pb.UserResp{
		UserId:   r.UserID,
		UserName: r.UserName,
		NickName: r.NickName,
	}
	return
}

func (u *UserSrv) UserRegister(ctx context.Context, req *pb.UserRegisterReq) (resp *pb.UserCommonResponse, err error) {

	zap.S().Info(logs.RunFuncName())

	resp = new(pb.UserCommonResponse)
	resp.Code = e2.SUCCESS
	err = dao.NewUserDao(ctx).CreateUser(req)
	if err != nil {
		resp.Code = e2.ERROR
		return
	}
	resp.Data = e2.GetMsg(int(resp.Code))
	return
}
