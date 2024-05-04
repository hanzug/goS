package main

import (
	"context"
	"fmt"
	"github.com/hanzug/goS/repository/redis"
	"go.uber.org/zap"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/hanzug/goS/app/search_engine/analyzer"
	"github.com/hanzug/goS/app/search_engine/repository/storage"
	"github.com/hanzug/goS/app/search_engine/service"
	"github.com/hanzug/goS/config"
	pb "github.com/hanzug/goS/idl/pb/search_engine"
	"github.com/hanzug/goS/loading"
	"github.com/hanzug/goS/pkg/discovery"
)

const SearchEngineService = "search_engine"

func main() {
	ctx := context.Background()
	loading.Loading()
	//rpc.Init()
	// bi_dao.InitDB() // TODO starrocks完善才开启
	analyzer.InitSeg()
	redis.InitRedis()
	storage.InitStorageDB(ctx)

	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress)
	grpcAddress := config.Conf.Services[SearchEngineService].Addr[0]
	defer etcdRegister.Stop()
	node := discovery.Server{
		Name: config.Conf.Domain[SearchEngineService].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定service
	pb.RegisterSearchEngineServiceServer(server, service.GetSearchEngineSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(node, 10); err != nil {
		panic(fmt.Sprintf("start service failed, err: %v", err))
	}
	logrus.Info("service started listen on ", grpcAddress)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
	defer zap.S().Sync()
}
