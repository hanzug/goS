package rpc

import (
	"context"
	"fmt"
	"github.com/hanzug/goS/config"
	"github.com/hanzug/goS/idl/pb/favorite"
	"github.com/hanzug/goS/idl/pb/index_platform"
	"github.com/hanzug/goS/idl/pb/search_engine"
	"github.com/hanzug/goS/idl/pb/user"
	"github.com/hanzug/goS/pkg/discovery"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

var (
	Register   *discovery.Resolver
	ctx        context.Context
	CancelFunc context.CancelFunc

	UserClient          user.UserServiceClient
	FavoriteClient      favorite.FavoritesServiceClient
	SearchEngineClient  search_engine.SearchEngineServiceClient
	IndexPlatformClient index_platform.IndexPlatformServiceClient
)

func Init() {
	zap.S().Info(logs.RunFuncName())

	// etcd解析器实例
	Register = discovery.NewResolver([]string{config.Conf.Etcd.Address})

	// 将解析器注册到grpc
	resolver.Register(Register)

	// 超时控制
	ctx, CancelFunc = context.WithTimeout(context.Background(), 3*time.Second)

	defer Register.Close()

	// 初始化微服务连接
	initClient(config.Conf.Domain["user"].Name, &UserClient)
	initClient(config.Conf.Domain["favorite"].Name, &FavoriteClient)
	initClient(config.Conf.Domain["search_engine"].Name, &SearchEngineClient)
	initClient(config.Conf.Domain["index_platform"].Name, &IndexPlatformClient)
}

func initClient(serviceName string, client interface{}) {

	zap.S().Info(logs.RunFuncName())

	// 连接对应的服务
	conn, err := connectServer(serviceName)
	zap.S().Info("connect ok", err, zap.Any("server name: ", serviceName))

	if err != nil {
		panic(err)
	}
	zap.S().Infof("Type of client: %T", client)

	switch c := client.(type) {
	case *user.UserServiceClient:
		*c = user.NewUserServiceClient(conn)
	case *favorite.FavoritesServiceClient:
		*c = favorite.NewFavoritesServiceClient(conn)
	case *search_engine.SearchEngineServiceClient:
		*c = search_engine.NewSearchEngineServiceClient(conn)
	case *index_platform.IndexPlatformServiceClient:
		*c = index_platform.NewIndexPlatformServiceClient(conn)
	default:
		panic("unsupported woker type")
	}
}

func connectServer(serviceName string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	addr := fmt.Sprintf("%s:///%s", Register.Scheme(), serviceName)
	zap.S().Info(addr)

	// Load balance
	if config.Conf.Services[serviceName].LoadBalance {
		log.Printf("load balance enabled for %s\n", serviceName)
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	}

	conn, err = grpc.DialContext(ctx, addr, opts...)
	return
}
