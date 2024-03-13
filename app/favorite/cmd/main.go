package main

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"

	"github.com/hanzug/goS/app/favorite/internal/service"
	"github.com/hanzug/goS/config"
	favoritePb "github.com/hanzug/goS/idl/pb/favorite"
	"github.com/hanzug/goS/loading"
	"github.com/hanzug/goS/pkg/discovery"
)

const ServerName = "favorite"

func main() {
	zap.S()
	loading.Loading()
	//rpc.Init()
	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress)
	grpcAddress := config.Conf.Services[ServerName].Addr[0]
	defer etcdRegister.Stop()
	node := discovery.Server{
		Name: config.Conf.Domain[ServerName].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定service
	favoritePb.RegisterFavoritesServiceServer(server, service.GetFavoriteSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(node, 10); err != nil {
		panic(fmt.Sprintf("start service failed, err: %v", err))
	}
	zap.S().Info("service started listen on ", grpcAddress)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}

	defer zap.S().Sync()
}
