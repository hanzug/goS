package loading

import (
	"github.com/hanzug/goS/config"
	logs "github.com/hanzug/goS/pkg/logger"
	"github.com/hanzug/goS/repository/mysql/db"
	"go.uber.org/zap"
)

func Loading() {

	zap.S().Info(logs.RunFuncName())

	// 读入配置
	config.InitConfig()

	// 初始化logger
	logs.InitLog()

	// 连接mysql
	db.InitDB()
}
