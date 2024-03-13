package main

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"

	"github.com/hanzug/goS/app/user/internal/service"
	"github.com/hanzug/goS/config"
	pb "github.com/hanzug/goS/idl/pb/user"
	"github.com/hanzug/goS/loading"
	"github.com/hanzug/goS/pkg/discovery"
	logs "github.com/hanzug/goS/pkg/logger"
)

const UserServiceName = "user"

func main() {

	zap.S().Info(logs.RunFuncName())

	// 初始化工作
	loading.Loading()

	// etcd地址
	etcdAddress := []string{config.Conf.Etcd.Address}

	// etcd 注册器
	etcdRegister := discovery.NewRegister(etcdAddress)
	defer etcdRegister.Stop()

	// user服务的grpc监听地址
	grpcAddress := config.Conf.Services[UserServiceName].Addr[0]

	// 服务节点信息
	userNode := discovery.Server{
		Name: config.Conf.Domain[UserServiceName].Name,
		Addr: grpcAddress,
	}

	// grpc初始化服务
	server := grpc.NewServer()
	defer server.Stop()

	// grpc初始化user服务
	pb.RegisterUserServiceServer(server, service.GetUserSrv())

	lis, err := net.Listen("tcp", grpcAddress)

	if err != nil {
		panic(err)
	}

	// 注册user服务节点
	if _, err := etcdRegister.Register(userNode, 10); err != nil {
		panic(fmt.Sprintf("start"))
	}
	zap.S().Info("service started listen on ", grpcAddress)

	// 开始服务套接字，直到发生错误
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
	defer zap.S().Sync()

}
