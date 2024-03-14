package recall

import (
	"context"
	"go.uber.org/zap"

	"github.com/hanzug/goS/types"
)

// SearchRecall 词条回归
func SearchRecall(ctx context.Context, query string) (res []*types.SearchItem, err error) {
	recallService := NewRecall()
	res, err = recallService.Search(ctx, query)
	if err != nil {
		zap.S().Errorf("SearchRecall-NewRecallServ:%+v", err)
		return
	}

	return
}

// SearchQuery 词条联想
func SearchQuery(query string) (res []string, err error) {
	recallService := NewRecall()
	res, err = recallService.SearchQuery(query)
	if err != nil {
		zap.S().Errorf("SearchRecall-NewRecallServ:%+v", err)
		return
	}

	return
}
