package dao

import (
	"context"
	"go.uber.org/zap"

	"gorm.io/gorm"

	favoritePb "github.com/hanzug/goS/idl/pb/favorite"
	"github.com/hanzug/goS/repository/mysql/db"
	"github.com/hanzug/goS/repository/mysql/model"
)

type FavoriteDao struct {
	*gorm.DB
}

func NewFavoriteDao(ctx context.Context) *FavoriteDao {
	return &FavoriteDao{db.NewDBClient(ctx)}
}

func (dao *FavoriteDao) ListFavorite(req *favoritePb.FavoriteListReq) (r []*model.Favorite, err error) {
	err = dao.DB.Model(&model.Favorite{}).
		Where("user_id = ?", req.UserId).Find(&r).Error
	if err != nil {
		return
	}
	return
}

func (dao *FavoriteDao) CreateFavorite(req *favoritePb.FavoriteCreateReq) (err error) {
	favorite := model.Favorite{
		FavoriteName: req.FavoriteName,
		UserID:       req.UserId,
	}
	if err = dao.DB.Create(&favorite).Error; err != nil {
		zap.S().Error("Insert Favorite Error:" + err.Error())
		return
	}

	return
}

func (dao *FavoriteDao) DeleteFavorite(req *favoritePb.FavoriteDeleteReq) (err error) {
	err = dao.DB.Where("favorite_id = ?", req.FavoriteId).
		Delete(model.Favorite{}).Error

	return
}

func (dao *FavoriteDao) UpdateFavorite(req *favoritePb.FavoriteUpdateReq) (err error) {
	fMap := make(map[string]interface{})
	fMap["favorite_name"] = req.FavoriteName
	err = dao.DB.Where("favorite_id = ?", req.FavoriteId).Updates(&fMap).Error
	if err != nil {
		return
	}

	return
}
