package service

import (
	"context"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"
	"sync"

	"github.com/hanzug/goS/app/search_engine/service/recall"
	"github.com/hanzug/goS/consts/e"
	pb "github.com/hanzug/goS/idl/pb/search_engine"
	"github.com/hanzug/goS/types"
)

var SearchEngineSrvIns *SearchEngineSrv
var SearchEngineSrvOnce sync.Once

type SearchEngineSrv struct {
	pb.UnimplementedSearchEngineServiceServer
}

func GetSearchEngineSrv() *SearchEngineSrv {
	SearchEngineSrvOnce.Do(func() {
		SearchEngineSrvIns = &SearchEngineSrv{}
	})
	return SearchEngineSrvIns
}

// SearchEngineSearch 搜索
func (s *SearchEngineSrv) SearchEngineSearch(ctx context.Context, req *pb.SearchEngineRequest) (resp *pb.SearchEngineResponse, err error) {

	zap.S().Info(logs.RunFuncName())

	resp = new(pb.SearchEngineResponse)
	resp.Code = e.SUCCESS
	query := req.Query
	sResult, err := recall.SearchRecall(ctx, query)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = err.Error()
		zap.S().Error("SearchEngineSearch-recall.SearchRecall", err)
		return
	}

	resp.SearchEngineInfoList, err = BuildSearchEngineResp(sResult)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = err.Error()
		zap.S().Error("SearchEngineSearch-BuildSearchEngineResp", err)
		return
	}
	resp.Count = int64(len(sResult))

	return
}

// WordAssociation 词语联想
func (s *SearchEngineSrv) WordAssociation(ctx context.Context, req *pb.SearchEngineRequest) (resp *pb.WordAssociationResponse, err error) {
	resp = new(pb.WordAssociationResponse)
	resp.Code = e.SUCCESS
	query := req.Query
	associationList, err := recall.SearchQuery(query)
	if err != nil {
		resp.Code = e.ERROR
		resp.Msg = err.Error()
		zap.S().Error("SearchEngineSearch-WordAssociation", err)
		return
	}
	resp.WordAssociationList = associationList

	return
}

func BuildSearchEngineResp(item []*types.SearchItem) (resp []*pb.SearchEngineList, err error) {
	resp = make([]*pb.SearchEngineList, 0)
	for _, v := range item {
		resp = append(resp, &pb.SearchEngineList{
			UrlId: v.DocId,
			Desc:  v.Content,
			Score: float32(v.Score),
		})
	}

	return
}
