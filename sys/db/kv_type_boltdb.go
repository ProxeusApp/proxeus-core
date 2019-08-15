package db

import (
	"os"
	"time"

	bolt "github.com/coreos/bbolt"
)

type KVTypeBoltDB struct {
	db *bolt.DB
}

var bucket = []byte("static_key")

func NewKVTypeBoldDB() (KVStoreIF, error) {
	return &KVTypeBoltDB{}, nil
}

func (me *KVTypeBoltDB) Create(path string) (err error) {
	me.db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	return
}

func (me *KVTypeBoltDB) Put(key *string, val []byte) error {
	return me.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}
		err = b.Put([]byte(*key), val)
		return err
	})
}

func (me *KVTypeBoltDB) Get(key *string) (resBytes []byte, err error) {
	err = me.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b != nil {
			resBytes = b.Get([]byte(*key))
		}
		return nil
	})
	return
}

func (me *KVTypeBoltDB) Delete(key *string) error {
	return me.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b != nil {
			return b.Delete([]byte(*key))
		}
		return nil
	})
}

func (me *KVTypeBoltDB) All() (keys []string, err error) {
	keys = make([]string, 0)
	err = me.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b != nil {
			return b.ForEach(func(k, v []byte) error {
				keys = append(keys, string(k))
				return nil
			})
		}
		return nil
	})
	return
}

func (me *KVTypeBoltDB) Close(delete bool) (err error) {
	p := me.db.Path()
	err = me.db.Close()
	if delete {
		err = os.Remove(p)
	}
	return
}
