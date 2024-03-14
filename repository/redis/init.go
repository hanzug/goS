package redis

import (
	"context"
	"fmt"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"

	"github.com/hanzug/goS/config"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client
var RedisContext = context.Background()

// InitRedis 在中间件中初始化redis链接
func InitRedis() {
	zap.S().Info(logs.RunFuncName())
	rConfig := config.Conf.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rConfig.RedisHost, rConfig.RedisPort),
		Username: rConfig.RedisUsername,
		Password: rConfig.RedisPassword,
		DB:       rConfig.RedisDbName,
	})
	_, err := client.Ping(RedisContext).Result()
	if err != nil {
		zap.S().Error(err)
		panic(err)
	}
	RedisClient = client
}
