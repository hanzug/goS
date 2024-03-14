package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hanzug/goS/app/search_engine/analyzer"
	"github.com/hanzug/goS/app/search_engine/service/recall"
	"github.com/hanzug/goS/config"
	log "github.com/hanzug/goS/pkg/logger"
	"github.com/hanzug/goS/repository/mysql/db"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	analyzer.InitSeg()
	log.InitLog()
	db.InitDB()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestRecall(t *testing.T) {
	q := "小岛"
	ctx := context.Background()
	searchItem, err := recall.SearchRecall(ctx, q)
	if err != nil {
		fmt.Println(err)
	}
	for i := range searchItem {
		fmt.Println(searchItem[i].Score, searchItem[i].DocId, searchItem[i].Content)
	}
}
