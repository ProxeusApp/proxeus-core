package db

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	bolt "go.etcd.io/bbolt"
)

type (
	KV struct {
		verbose       bool
		diskStore     *bolt.DB
		diskStorePath string
		diskStoreDir  string
		diskStoreName string
	}
)

var (
	ErrNotFound = errors.New("key does not exist")
)

func NewKV(path string) (*KV, error) {
	pc := &KV{
		diskStorePath: path,
	}
	pc.checkPath()
	err := pc.openBoltDB()
	if err != nil {
		return nil, err
	}
	return pc, nil
}

func (me *KV) DiskPath() string {
	return me.diskStorePath
}

func (me *KV) DiskFilename() string {
	return me.diskStoreName
}

func (me *KV) DiskDir() string {
	return me.diskStoreDir
}

func (me *KV) Get(bucket, key string, ref interface{}) error {
	return me.diskStore.View(func(tx *bolt.Tx) error {
		b2 := tx.Bucket([]byte(bucket))
		if b2 != nil {
			err := json.Unmarshal(b2.Get([]byte(key)), ref)
			if err != nil {
				return err
			}
			return nil
		}
		return ErrNotFound
	})
}

func (me *KV) Remove(bucket, key string) error {
	return me.diskStore.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return b.Delete([]byte(key))
	})
}

func (me *KV) Iterate(bucket string, cb func(key, value []byte) error) error {
	return me.diskStore.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			c := b.Cursor()
			var err error
			for k, v := c.First(); k != nil; k, v = c.Next() {
				err = cb(k, v)
				if err != nil {
					return err
				}
			}
		} else {
			return ErrNotFound
		}
		return nil
	})
}

func (me *KV) Put(bucket, key string, val interface{}) error {
	dataBts, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return me.diskStore.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		err = b.Put([]byte(key), dataBts)
		return err
	})
}

func (me *KV) openBoltDB() (err error) {
	me.diskStore, err = bolt.Open(me.diskStorePath, 0600, &bolt.Options{NoGrowSync: true, Timeout: 1 * time.Second})
	return
}

func (me *KV) checkPath() {
	me.diskStorePath = strings.TrimSpace(me.diskStorePath)
	me.diskStoreName = "."
	if !strings.HasSuffix(me.diskStorePath, string(os.PathSeparator)) {
		me.diskStoreName = filepath.Base(me.diskStorePath)
	}
	if me.diskStoreName == "." {
		me.diskStoreName = "kv"
	}
	me.diskStoreDir = filepath.Dir(me.diskStorePath)
	me.diskStorePath = filepath.Join(me.diskStoreDir, me.diskStoreName)
}

func (me *KV) ensureDir() error {
	var err error
	_, err = os.Stat(me.diskStoreDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(me.diskStoreDir, 0750)
		if err != nil {
			return err
		}
	}
	return nil
}

func (me *KV) pathExists() bool {
	_, err := os.Stat(me.diskStorePath)
	return !os.IsNotExist(err)
}

func (me *KV) Close() {
	me.diskStore.Close()
}
