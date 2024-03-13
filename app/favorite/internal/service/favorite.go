package service

import (
	"context"
	"sync"

	"github.com/hanzug/goS/app/favorite/internal/repository/db/dao"
	"github.com/hanzug/goS/consts/e"
	pb "github.com/hanzug/goS/idl/pb/favorite"
	"github.com/hanzug/goS/repository/mysql/model"
)

var FavoriteSrvIns *FavoriteSrv
var FavoriteSrvOnce sync.Once

type FavoriteSrv struct {
	pb.UnimplementedFavoritesServiceServer
}

// GetFavoriteSrv 返回grpc结构体
func GetFavoriteSrv() *FavoriteSrv {
	FavoriteSrvOnce.Do(func() {
		FavoriteSrvIns = &FavoriteSrv{}
	})
	return FavoriteSrvIns
}

// FavoriteCreate 创建收藏夹
func (s *FavoriteSrv) FavoriteCreate(ctx context.Context, req *pb.FavoriteCreateReq) (resp *pb.FavoriteCommonResponse, err error) {
	resp = new(pb.FavoriteCommonResponse)
	resp.Code = e.SUCCESS
	err = dao.NewFavoriteDao(ctx).CreateFavorite(req)
	if err != nil {
		resp.Error = err.Error()
		return
	}

	resp.Msg = e.GetMsg(int(resp.Code))
	return
}

func (s *FavoriteSrv) FavoriteList(ctx context.Context, req *pb.FavoriteListReq) (resp *pb.FavoriteListResponse, err error) {
	resp = new(pb.FavoriteListResponse)
	f, err := dao.NewFavoriteDao(ctx).ListFavorite(req)
	resp.Code = e.SUCCESS
	if err != nil {
		resp.Code = e.ERROR
		return
	}
	for i := range f {
		resp.Items = append(resp.Items, &pb.FavoriteListItemResp{
			FavoriteId:   f[i].FavoriteID,
			FavoriteName: f[i].FavoriteName,
		})
	}

	return
}

func (s *FavoriteSrv) FavoriteUpdate(ctx context.Context, req *pb.FavoriteUpdateReq) (resp *pb.FavoriteCommonResponse, err error) {
	resp = new(pb.FavoriteCommonResponse)
	resp.Code = e.SUCCESS
	err = dao.NewFavoriteDao(ctx).UpdateFavorite(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Error = err.Error()
		return
	}

	resp.Msg = e.GetMsg(int(resp.Code))
	return resp, nil
}

func (s *FavoriteSrv) FavoriteDelete(ctx context.Context, req *pb.FavoriteDeleteReq) (resp *pb.FavoriteCommonResponse, err error) {
	resp = new(pb.FavoriteCommonResponse)
	resp.Code = e.SUCCESS
	err = dao.NewFavoriteDao(ctx).DeleteFavorite(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Error = err.Error()
		return
	}

	resp.Msg = e.GetMsg(int(resp.Code))
	return
}

func (s *FavoriteSrv) FavoriteDetailCreate(ctx context.Context, req *pb.FavoriteDetailCreateReq) (resp *pb.FavoriteCommonResponse, err error) {
	resp = new(pb.FavoriteCommonResponse)
	resp.Code = e.SUCCESS
	err = dao.NewFavoriteDetailDao(ctx).CreateFavoriteDetail(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Error = err.Error()
		return
	}
	resp.Msg = e.GetMsg(int(resp.Code))
	return
}

func (s *FavoriteSrv) FavoriteDetailDelete(ctx context.Context, req *pb.FavoriteDetailDeleteReq) (resp *pb.FavoriteCommonResponse, err error) {
	resp = new(pb.FavoriteCommonResponse)
	resp.Code = e.SUCCESS
	err = dao.NewFavoriteDetailDao(ctx).DeleteFavoriteDetail(req)
	if err != nil {
		resp.Code = e.ERROR
		resp.Error = err.Error()
		return
	}

	resp.Msg = e.GetMsg(int(resp.Code))
	return
}

func (s *FavoriteSrv) FavoriteDetailList(ctx context.Context, req *pb.FavoriteDetailListReq) (resp *pb.FavoriteDetailListResponse, err error) {
	resp = new(pb.FavoriteDetailListResponse)
	resp.Code = e.SUCCESS
	fdResp, err := dao.NewFavoriteDetailDao(ctx).ListFavoriteDetail(req)
	if err != nil {
		resp.Code = e.ERROR
		return
	}

	resp.Items = BuildFavoriteDetails(fdResp)
	return
}

func BuildFavoriteDetails(item []*model.Favorite) (fList []*pb.FavoriteResp) {
	for _, v := range item {
		f := BuildFavoriteDetail(v)
		fList = append(fList, f)
	}
	return fList
}

func BuildFavoriteDetail(item *model.Favorite) *pb.FavoriteResp {
	return &pb.FavoriteResp{
		FavoriteId:   item.FavoriteID,
		FavoriteName: item.FavoriteName,
		UserId:       item.UserID,
		UrlInfo:      BuildUrlInfos(item.FavoriteDetail),
	}
}

func BuildUrlInfo(item *model.FavoriteDetail) *pb.UrlModel {
	return &pb.UrlModel{
		UrlId: item.UrlID,
		Url:   item.Url,
		Desc:  item.Desc,
	}
}

func BuildUrlInfos(item []*model.FavoriteDetail) (urlList []*pb.UrlModel) {
	for _, v := range item {
		u := BuildUrlInfo(v)
		urlList = append(urlList, u)
	}
	return urlList
}
