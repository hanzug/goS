package main

import (
	"context"
	"fmt"
	"github.com/hanzug/goS/app/index_platform/analyzer"
	"github.com/hanzug/goS/app/index_platform/cmd/kfk_register"
	"github.com/hanzug/goS/loading"
	"github.com/hanzug/goS/pkg/kafka"
	"github.com/hanzug/goS/repository/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"

	"github.com/hanzug/goS/app/index_platform/service"
	"github.com/hanzug/goS/config"
	"github.com/hanzug/goS/idl/pb/index_platform"
	"github.com/hanzug/goS/pkg/discovery"
	logs "github.com/hanzug/goS/pkg/logger"
)

const (
	IndexPlatformServerName = "index_platform"
)

func main() {
	zap.S().Info(logs.RunFuncName())

	ctx := context.Background()
	// 加载配置
	loading.Loading()

	redis.InitRedis()

	// 分词器初始化
	analyzer.InitSeg()

	// 连接kafka
	kafka.InitKafka()

	// 启动kafka消费者
	kfk_register.RegisterJob(ctx)

	// 注册服务
	_ = registerIndexPlatform()
}

// registerIndexPlatform 注册索引平台服务
func registerIndexPlatform() (err error) {
	etcdAddress := []string{config.Conf.Etcd.Address}
	etcdRegister := discovery.NewRegister(etcdAddress)
	defer etcdRegister.Stop()

	grpcAddress := config.Conf.Services[IndexPlatformServerName].Addr[0]

	node := discovery.Server{
		Name: config.Conf.Domain[IndexPlatformServerName].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()

	index_platform.RegisterIndexPlatformServiceServer(server, service.GetIndexPlatformSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err = etcdRegister.Register(node, 10); err != nil {
		panic(fmt.Sprintf("start service failed, err: %v", err))
	}
	zap.S().Info("service started listen on ", grpcAddress)
	if err = server.Serve(lis); err != nil {
		panic(err)
	}
	defer zap.S().Sync()
	return
}
