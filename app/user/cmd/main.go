package main

import (
	"fmt"
	"github.com/hanzug/goS/config"
	pb "github.com/hanzug/goS/idl/pb/user"
	"github.com/hanzug/goS/loading"
	"github.com/hanzug/goS/pkg/discovery"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

const UserServiceName = "user"

func main() {

	zap.S().Info(logs.RunFuncName())

	loading.Loading()
	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress)

	//grpc 地址
	grpcAddress := config.Conf.Services[UserServiceName].Addr[0]
	defer etcdRegister.Stop()
	userNode := discovery.Server{
		Name: config.Conf.Domain[UserServiceName].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定service
	pb.RegisterUserServiceServer(server, service.GetUserSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	//把userService注册到etcd，使用grpc来交互。
	if _, err := etcdRegister.Register(userNode, 10); err != nil {
		panic(fmt.Sprintf("start service failed, err: %v", err))
	}
	zap.S().Info("service started listen on ", grpcAddress)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
	defer zap.S().Sync()
}
