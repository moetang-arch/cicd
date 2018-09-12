package core

import (
	"time"

	"github.com/boltdb/bolt"
)

type StoreService struct {
	db             *bolt.DB
	initHelperImpl initHelper
}

func (this *StoreService) GetHelper() initHelper {
	return this.initHelperImpl
}

type HelperFactory interface {
	GetHelper() initHelper
}

type initHelper struct {
	store *StoreService
}

func (this initHelper) CreateTable(tableName string) error {
	return this.store.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(tableName))
		return err
	})
}

func NewStoreService(path string) (*StoreService, error) {
	option := new(bolt.Options)
	option.Timeout = 10 * time.Second
	db, err := bolt.Open(path, 0644, option)
	if err != nil {
		return nil, err
	}

	ss := new(StoreService)
	ss.db = db
	ss.initHelperImpl.store = ss

	return ss, nil
}

func (this *StoreService) Close() error {
	return this.db.Close()
}

func (this *StoreService) Init(fn func(helperFactory HelperFactory) error) error {
	return fn(this)
}
