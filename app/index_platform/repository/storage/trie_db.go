package storage

import (
	logs "github.com/hanzug/goS/pkg/logger"
	"go.uber.org/zap"
	"os"

	bolt "go.etcd.io/bbolt"

	"github.com/hanzug/goS/consts"
	"github.com/hanzug/goS/pkg/trie"
)

type TrieDB struct {
	file *os.File
	db   *bolt.DB
}

// NewTrieDB 初始化trie
func NewTrieDB(filePath string) *TrieDB { // TODO: 先都放在一个下面吧，后面再lb到多个文件

	zap.S().Info(logs.RunFuncName())

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		zap.S().Error(err)
	}

	db, err := bolt.Open(filePath, 0600, nil)
	if err != nil {
		zap.S().Error(err)
		return nil
	}

	return &TrieDB{f, db}
}

func (d *TrieDB) StorageDict(trieTree *trie.Trie) (err error) {

	zap.S().Info(logs.RunFuncName())

	trieByte, _ := trieTree.Root.Children.MarshalJSON()
	err = d.PutTrieTree([]byte(consts.TrieTreeBucket), trieByte)

	return
}

// GetTrieTreeInfo 获取 trie tree
func (d *TrieDB) GetTrieTreeInfo() (trieTree *trie.Trie, err error) {

	zap.S().Info(logs.RunFuncName())

	v, err := d.GetTrieTree([]byte(consts.TrieTreeBucket))
	if err != nil {
		return
	}

	trieTree = trie.NewTrie()
	err = trieTree.UnmarshalJSON(v)

	return
}

// PutTrieTree 存储
func (d *TrieDB) PutTrieTree(key, value []byte) error {

	zap.S().Info(logs.RunFuncName())

	return Put(d.db, consts.TrieTreeBucket, key, value)
}

// GetTrieTree 通过term获取value
func (d *TrieDB) GetTrieTree(key []byte) (value []byte, err error) {

	zap.S().Info(logs.RunFuncName())

	return Get(d.db, consts.TrieTreeBucket, key)
}

// Close 关闭db
func (d *TrieDB) Close() error {

	zap.S().Info(logs.RunFuncName())

	return d.db.Close()
}
