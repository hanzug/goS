package storage

import (
	logs "github.com/hanzug/goS/pkg/logger"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

// Put 通过bolt写入数据
func Put(db *bolt.DB, bucket string, key []byte, value []byte) error {

	zap.S().Info(logs.RunFuncName())

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})
}

// Get 通过bolt获取数据
func Get(db *bolt.DB, bucket string, key []byte) (r []byte, err error) {

	zap.S().Info(logs.RunFuncName())

	err = db.View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(bucket))
		r = b.Get(key)
		if r == nil { // 如果是空的话，直接创建这个key，然后返回这个key的初始值，也就是0
			r = []byte("0")
			return
		}
		return
	})

	return
}
