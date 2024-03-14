package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/hanzug/goS/consts"
	"github.com/hanzug/goS/repository/mysql/db"
	"github.com/hanzug/goS/repository/mysql/model"
	"github.com/hanzug/goS/types"
)

type InputDataDao struct {
	*gorm.DB
}

func NewInputDataDao(ctx context.Context) *InputDataDao {
	return &InputDataDao{db.NewDBClient(ctx)}
}

func (d *InputDataDao) CreateInputData(in *model.InputData) (err error) {
	return d.DB.Model(&model.InputData{}).Create(&in).Error
}

func (d *InputDataDao) BatchCreateInputData(in []*model.InputData) (err error) {
	return d.DB.Model(&model.InputData{}).CreateInBatches(&in, consts.BatchCreateSize).Error
}

func (d *InputDataDao) ListInputData() (in []*model.InputData, err error) {
	err = d.DB.Model(&model.InputData{}).Where("is_index = ?", false).
		Find(&in).Error

	return
}

// ListInputDataByDocIds 根据传进来的 doc id 获取所有的信息
func (d *InputDataDao) ListInputDataByDocIds(docIds []uint32) (in []*types.SearchItem, err error) {
	err = d.DB.Model(&model.InputData{}).
		Where("doc_id IN ?", docIds).
		Select("doc_id," +
			"title," +
			"body AS content," +
			"url," +
			"score AS content_score").
		Find(&in).Error

	return
}

func (d *InputDataDao) UpdateInputDataByIds(ids []int64) (err error) {
	err = d.DB.Model(&model.InputData{}).Where("id IN ?", ids).
		Update("is_index", true).Error

	return
}
