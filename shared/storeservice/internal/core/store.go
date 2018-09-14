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

func (this *StoreService) InsertInto(tableName string, key []byte, data []byte) error {
	return this.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		err := b.Put(key, data)
		if err != nil {
			return err
		}
		return nil
	})
}

func (this *StoreService) Get(tableName string, key []byte) ([]byte, error) {
	var result []byte
	err := this.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		result = b.Get(key)
		return nil
	})
	return result, err
}

func (this *StoreService) Modify(tableName string, key []byte, data []byte) (modified bool, err error) {
	err = this.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tableName))
		v := b.Get(key)
		if v == nil {
			modified = false
			return nil
		}
		e := b.Put(key, data)
		if e != nil {
			return e
		}
		return nil
	})
	return
}
