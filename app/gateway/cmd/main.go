package main

import (
	"fmt"
	"github.com/hanzug/goS/app/gateway/rpc"
	"github.com/hanzug/goS/config"
	"github.com/hanzug/goS/loading"
	"github.com/hanzug/goS/pkg/discovery"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
	"net/http"
	"time"
)

func main() {
	zap.S().Info(logs.RunFuncName())

	// 加载配置文件和初始化日志等
	loading.Loading()
	// 初始化RPC服务
	rpc.Init()
	etcdAddress := []string{config.Conf.Etcd.Address}
	etcdRegister := discovery.NewResolver(etcdAddress)
	zap.S().Info("loading ok")
	zap.S().Sync()

	defer etcdRegister.Close()
	//注册etcd解析器到grpc中
	resolver.Register(etcdRegister)
	zap.S().Sync()
	go startListen() // 转载路由
	// {
	// 	osSignals := make(chan os.Signal, 1)
	// 	signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	// 	s := <-osSignals
	// 	fmt.Println("exit! ", s)
	// }
	zap.S().Sync()
}

func startListen() {
	zap.S().Info(logs.RunFuncName())
	ginRouter := routes.NewRouter()
	server := &http.Server{
		Addr:           config.Conf.Server.Port,
		Handler:        ginRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("绑定HTTP到 %s 失败！可能是端口已经被占用，或用户权限不足 \n", config.Conf.Server.Port)
		fmt.Println(err)
		return
	}
	fmt.Printf("gateway listen on :%v \n", config.Conf.Server.Port)
	// go func() {
	// 	// TODO 优雅关闭 有点问题，后续优化一下
	// 	shutdown.GracefullyShutdown(service)
	// }()
	defer zap.S().Sync()
}
