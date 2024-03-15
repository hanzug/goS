package storage

import (
	"context"
	"errors"
	//"fmt"
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"
	"os"

	bolt "go.etcd.io/bbolt"

	"github.com/hanzug/goS/consts"
	"github.com/hanzug/goS/repository/redis"
)

type KvInfo struct {
	Key   []byte
	Value []byte
}

var GlobalInvertedDB []*InvertedDB

type InvertedDB struct {
	file   *os.File
	db     *bolt.DB
	offset int64
}

// InitInvertedDB 初始化倒排索引库
func InitInvertedDB(ctx context.Context) []*InvertedDB {
	zap.S().Info(logs.RunFuncName())

	dbs := make([]*InvertedDB, 0)
	filePath, _ := redis.ListInvertedPath(ctx, redis.InvertedIndexDbPathKey)
	i := 0
	for _, file := range filePath {
		i++
		if i == 2 {
			break
		}
		f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			zap.S().Error(err)
		}
		stat, err := f.Stat()
		if err != nil {
			zap.S().Error(err)
		}
		db, err := bolt.Open(file, 0600, nil)
		if err != nil {
			zap.S().Error(err)
		}
		zap.S().Info(stat.Size())

		dbs = append(dbs, &InvertedDB{f, db, stat.Size()})
	}
	if len(filePath) == 0 {
		panic(errors.New("没有索引库...请先创建索引库"))
	}
	GlobalInvertedDB = dbs
	zap.S().Info(logs.RunFuncName())

	return nil
}

// NewInvertedDB 新建一个inverted
func NewInvertedDB(termName, postingsName string) *InvertedDB {
	zap.S().Info(logs.RunFuncName())

	f, err := os.OpenFile(postingsName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		zap.S().Error(err)
	}
	stat, err := f.Stat()
	if err != nil {
		zap.S().Error(err)
	}
	zap.S().Infof("start op bolt:%v", termName)
	db, err := bolt.Open(termName, 0600, nil)
	if err != nil {
		zap.S().Error(err)
	}
	return &InvertedDB{f, db, stat.Size()}
}

// GetInverted 通过term获取value
func (t *InvertedDB) GetInverted(key []byte) (value []byte, err error) {
	return Get(t.db, consts.InvertedBucket, key)
}

//// GetInvertedDoc 根据地址获取读取文件
//func (t *InvertedDB) GetInvertedDoc(offset int64, size int64) ([]byte, error) {
//	zap.S().Info(logs.RunFuncName())
//	page := os.Getpagesize()
//	b, err := Mmap(int(t.file.Fd()), offset/int64(page), int(offset+size))
//	if err != nil {
//		return nil, fmt.Errorf("GetDocinfo Mmap err: %v", err)
//	}
//	return b[offset : offset+size], nil
//}

func (t *InvertedDB) Close() {
	t.file.Close()
	t.db.Close()
}
