package storage

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"

	"github.com/RoaringBitmap/roaring"
	bolt "go.etcd.io/bbolt"

	"github.com/hanzug/goS/consts"
	"github.com/hanzug/goS/types"
)

type KvInfo struct {
	Key   []byte
	Value []byte
}

type InvertedDB struct { // 后续做mmap(这个好难)
	file   *os.File
	db     *bolt.DB
	offset int64
}

// NewInvertedDB 新建一个inverted
func NewInvertedDB(invertedName string) *InvertedDB {
	f, err := os.OpenFile(invertedName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		zap.S().Error(err)
	}
	stat, err := f.Stat()
	if err != nil {
		zap.S().Error(err)
	}
	db, err := bolt.Open(invertedName, 0600, nil)
	if err != nil {
		zap.S().Error(err)
	}

	return &InvertedDB{f, db, stat.Size()}
}

// StoragePostings 存储 倒排索引表
func (t *InvertedDB) StoragePostings(token string, values []byte) (err error) {
	return t.PutInverted([]byte(token), values)
}

// PutInverted 插入term
func (t *InvertedDB) PutInverted(key, value []byte) error {
	return Put(t.db, consts.InvertedBucket, key, value)
}

// GetInverted 通过term获取value
func (t *InvertedDB) GetInverted(key []byte) (value []byte, err error) {
	return Get(t.db, consts.InvertedBucket, key)
}

// GetInvertedInfo 获取倒排地址
func (t *InvertedDB) GetInvertedInfo(token string) (p *types.InvertedInfo, err error) {
	c, err := t.GetInverted([]byte(token))
	if err != nil {
		return
	}

	if len(c) == 0 {
		err = errors.New("暂无此token")
		return
	}
	output := roaring.New()
	_ = output.UnmarshalBinary(c)
	p = &types.InvertedInfo{
		Token:  token,
		DocIds: output,
	}
	return
}

// GetInvertedDoc 根据地址获取读取文件
func (t *InvertedDB) GetInvertedDoc(offset int64, size int64) ([]byte, error) {
	page := os.Getpagesize()
	b, err := Mmap(int(t.file.Fd()), offset/int64(page), int(offset+size))
	if err != nil {
		return nil, fmt.Errorf("GetDocinfo Mmap err: %v", err)
	}
	return b[offset : offset+size], nil
}

func (t *InvertedDB) Close() {
	err := t.file.Close()
	if err != nil {
		return
	}
	err = t.db.Close()
	if err != nil {
		return
	}
}
