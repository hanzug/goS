package dao

import (
	"context"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"

	"gorm.io/gorm"

	"github.com/hanzug/goS/consts"
	db "github.com/hanzug/goS/repository/mysql/db"
	"github.com/hanzug/goS/repository/mysql/model"
)

type InputDataDao struct {
	*gorm.DB
}

func NewInputDataDao(ctx context.Context) *InputDataDao {
	zap.S().Info(logs.RunFuncName())
	return &InputDataDao{db.NewDBClient(ctx)}
}

func (d *InputDataDao) CreateInputData(in *model.InputData) (err error) {
	zap.S().Info(logs.RunFuncName())
	return d.DB.Model(&model.InputData{}).Create(&in).Error
}

func (d *InputDataDao) BatchCreateInputData(in []*model.InputData) (err error) {
	zap.S().Info(logs.RunFuncName())
	return d.DB.Model(&model.InputData{}).CreateInBatches(&in, consts.BatchCreateSize).Error
}

func (d *InputDataDao) ListInputData() (in []*model.InputData, err error) {
	zap.S().Info(logs.RunFuncName())
	err = d.DB.Model(&model.InputData{}).Where("is_index = ?", false).
		Find(&in).Error

	return
}

func (d *InputDataDao) UpdateInputDataByIds(ids []int64) (err error) {
	zap.S().Info(logs.RunFuncName())
	err = d.DB.Model(&model.InputData{}).Where("id IN ?", ids).
		Update("is_index", true).Error

	return
}
