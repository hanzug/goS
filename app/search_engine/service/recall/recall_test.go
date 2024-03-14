package recall

import (
	"context"
	"fmt"
	"testing"

	"github.com/hanzug/goS/app/search_engine/repository/storage"
	"github.com/hanzug/goS/config"
	log "github.com/hanzug/goS/pkg/logger"
	"github.com/hanzug/goS/repository/redis"
)

func TestMain(m *testing.M) {
	// 这个文件相对于config.yaml的位置
	re := config.ConfigReader{FileName: "../../../../config/config.yaml"}
	config.InitConfigForTest(&re)
	log.InitLog()
	redis.InitRedis()
	fmt.Println("Write tests on values: ", config.Conf)
	m.Run()
}

func TestGetTrieTreeFromRedis(t *testing.T) {
	ctx := context.Background()
	storage.InitGlobalTrieDB(ctx)
	for _, v := range storage.GlobalTrieDB {
		tree, err := v.GetTrieTreeDict()
		if err != nil {
			fmt.Println("tree ", err)
		}
		tree.TraverseForRecall()
	}

}
