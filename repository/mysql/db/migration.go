package db

import (
	"go.uber.org/zap"
	"os"

	"github.com/hanzug/goS/repository/mysql/model"
)

func migration() {
	// 自动迁移模式
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.User{},
			&model.InputData{},
			&model.Favorite{},
			&model.FavoriteDetail{},
		)
	if err != nil {
		zap.S().Infof("register table fail")
		os.Exit(0)
	}
	zap.S().Infof("register table success")
}
