package db

import (
	"github.com/hanzug/goS/config"
	"go.uber.org/zap"

	logs "github.com/hanzug/goS/pkg/logger"

	"gorm.io/gorm"
)

var _db *gorm.DB

func InitDB() {
	zap.S().Info(logs.RunFuncName())

	mConfig := config.Conf.MySQL

	host :=
}
