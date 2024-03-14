package woker

import (
	"context"
	"fmt"
	"testing"

	"github.com/hanzug/goS/app/index_platform/analyzer"
	"github.com/hanzug/goS/app/index_platform/repository/db/dao"
	"github.com/hanzug/goS/app/index_platform/trie"
	"github.com/hanzug/goS/app/mapreduce/mr/input_data_mr"
	"github.com/hanzug/goS/app/mapreduce/rpc"
	"github.com/hanzug/goS/config"
	"github.com/hanzug/goS/pkg/kfk"
	"github.com/hanzug/goS/repository/mysql/db"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	db.InitDB()
	trie.InitTrieTree()
	analyzer.InitSeg()
	rpc.Init()
	kfk.InitKafka()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestWorker(t *testing.T) {
	ctx := context.Background()
	dao.InitMysqlDirectUpload(ctx)
	Worker(ctx, input_data_mr.Map, input_data_mr.Reduce)
}
