package rpc

import (
	"context"
	"errors"

	"github.com/hanzug/goS/consts/e"
	pb "github.com/hanzug/goS/idl/pb/search_engine"
)

func SearchEngineSearch(ctx context.Context, req *pb.SearchEngineRequest) (r *pb.SearchEngineResponse, err error) {
	r, err = SearchEngineClient.SearchEngineSearch(ctx, req)
	if err != nil {
		return
	}

	if r.Code != e.SUCCESS {
		err = errors.New(r.Msg)
		return
	}

	return
}

func WordAssociation(ctx context.Context, req *pb.SearchEngineRequest) (r *pb.WordAssociationResponse, err error) {
	r, err = SearchEngineClient.WordAssociation(ctx, req)
	if err != nil {
		return
	}

	if r.Code != e.SUCCESS {
		err = errors.New(r.Msg)
		return
	}

	return
}
