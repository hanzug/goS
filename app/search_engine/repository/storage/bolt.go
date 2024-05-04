package storage

import (
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

// Put 通过bolt写入数据
func Put(db *bolt.DB, bucket string, key []byte, value []byte) error {
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
	err = db.Update(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			zap.S().Info("b is nil", zap.Any("bucket", bucket))
			b, err = tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				zap.S().Error("create bucket failed", zap.Any("bucket", bucket))
				return err
			}
			zap.S().Info("new bucket created", zap.Any("new b", b))
		}
		r = b.Get(key)
		if r == nil {
			r = []byte("0")
		}
		return
	})

	return
}
