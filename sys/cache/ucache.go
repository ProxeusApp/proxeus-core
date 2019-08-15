package cache

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	DiskCacheSpeedMode       StoreType = 0 // DEFAULT speed mode handles persistence to disk in the background in batches - read from disk write fast uses mem only until persisted
	MemCache                 StoreType = 1 // mem cache only - read fast write fast
	DiskCache                StoreType = 2 // disk cache only will persist right away - less memory usage but slow - read slow write slow
	MemAndDiskCache          StoreType = 3 // keep cache in memory and persist it right away - read fast write slow
	MemAndDiskCacheSpeedMode StoreType = 4 // speed mode handles persistence to disk in the background in batches - read fast write fast but uses more memory

	None                  ValueBeType = 0 // DEFAULT don't do anything with the value
	CallOnExpire          ValueBeType = 1 // call OnExpire() when it is being removed
	CallOnLoadAndOnExpire ValueBeType = 2 // call OnLoad() when the value was loaded from the disk and call OnExpire() when it is being removed

	//default reflection value behaviour method names
	onExpire = "OnExpire"
	onLoad   = "OnLoad"
)

type (
	//CallBack method to get life cycle feedback on the cache entries
	//In usage of OnExpireFunc there is an exception. If the value is not loaded on memory and it is expiring "val interface{} will be nil"!
	CallBack       func(key string, val interface{})
	StoreType      int
	ValueBeType    int
	ValueBehaviour struct {
		ValueBeType        ValueBeType //default None
		OnLoadMethodName   string      //OnLoad method for using reflection on values. Leave empty if OnLoadFunc is set for not using reflection.
		OnExpireMethodName string      //OnExpire method for using reflection on values. Leave empty if OnExpireFunc is set for not using reflection.
		OnLoadFunc         CallBack    //Global OnLoad callback to prevent from reflection.
		OnExpireFunc       CallBack    //Global OnExpire callback to prevent from reflection. If the value is not loaded on memory and it is expiring "val interface{} will be nil"!
	}
	UCacheConfig struct {
		ValueBehaviour *ValueBehaviour // if you have complex objects you want to persist, you can handle the lifecycle of them with ValueBehaviour
		Verbose        bool            // print out log
		DiskStorePath  string          // provide a directory or a file - if filename is not provided it will be named cache
		DefaultExpiry  time.Duration   // default expiry for calling Put - PutWithOtherExpiry makes it possible to set another one
		StoreType      StoreType       // store type has a huge effect on memory usage and performance - choose use case specific
		ExtendExpiry   bool            // extend expiry when access cache over Get
		NoExpiry       bool            // disable garbage collection
	}

	value struct {
		m                  sync.RWMutex
		access             time.Time
		duration           time.Duration
		persisted          bool // used for speed mode
		loadingFromDisk    bool
		loadingFromDiskErr error
		val                interface{}
	}

	UCache struct {
		noExpiry       bool
		verbose        bool
		valueBehaviour *ValueBehaviour
		storeType      StoreType
		memCache       bool
		diskCache      bool
		memStore       map[string]*value

		diskSyncStopSignal chan bool
		diskSyncStopped    bool
		diskSyncLock       sync.Mutex
		diskStore          *bolt.DB
		diskStorePath      string
		diskStoreDir       string
		diskStoreName      string

		gcStoppedLock  sync.Mutex
		gcStopped      bool
		gcStop         chan bool
		diskSyncNotify chan bool
		expiry         time.Duration

		memLock             sync.RWMutex
		defaultExtendExpiry bool
		speedMode           bool
		lastPut             time.Time
		lastGCRun           time.Time
		nextToBeCleanedUp   [][]byte
	}

	diskMeta struct {
		A time.Time     //Access
		D time.Duration //Duration
	}
)

var (
	ErrNotSupportedWithoutMemCache = errors.New("not supported without mem cache")
	ErrNotFound                    = errors.New("key does not exist")
	bucketMeta                     = []byte("meta")
	bucketData                     = []byte("data")
)

func NewUCache(uCacheConfig UCacheConfig) (*UCache, error) {
	if !uCacheConfig.NoExpiry && uCacheConfig.DefaultExpiry == 0 {
		return nil, fmt.Errorf("default expiry must be set higher than 0")
	}
	pc := &UCache{
		expiry:              uCacheConfig.DefaultExpiry,
		memCache:            uCacheConfig.StoreType == MemCache || uCacheConfig.StoreType == MemAndDiskCache || uCacheConfig.StoreType == DiskCacheSpeedMode || uCacheConfig.StoreType == MemAndDiskCacheSpeedMode,
		diskCache:           uCacheConfig.StoreType == DiskCache || uCacheConfig.StoreType == MemAndDiskCache || uCacheConfig.StoreType == DiskCacheSpeedMode || uCacheConfig.StoreType == MemAndDiskCacheSpeedMode,
		speedMode:           uCacheConfig.StoreType == DiskCacheSpeedMode || uCacheConfig.StoreType == MemAndDiskCacheSpeedMode,
		verbose:             uCacheConfig.Verbose,
		diskStorePath:       uCacheConfig.DiskStorePath,
		storeType:           uCacheConfig.StoreType,
		noExpiry:            uCacheConfig.NoExpiry,
		gcStopped:           true,
		memStore:            make(map[string]*value),
		nextToBeCleanedUp:   make([][]byte, 0),
		defaultExtendExpiry: uCacheConfig.ExtendExpiry,
	}
	if uCacheConfig.ValueBehaviour == nil {
		pc.valueBehaviour = &ValueBehaviour{ValueBeType: None}
	} else {
		pc.valueBehaviour = &ValueBehaviour{
			ValueBeType:        uCacheConfig.ValueBehaviour.ValueBeType,
			OnLoadMethodName:   uCacheConfig.ValueBehaviour.OnLoadMethodName,
			OnExpireMethodName: uCacheConfig.ValueBehaviour.OnExpireMethodName,
			OnLoadFunc:         uCacheConfig.ValueBehaviour.OnLoadFunc,
			OnExpireFunc:       uCacheConfig.ValueBehaviour.OnExpireFunc,
		}
		pc.valueBehaviour.check()
	}

	if pc.diskCache {
		pc.checkPath()
		err := pc.openBoltDB()
		if err != nil {
			return nil, err
		}
		pc.diskSyncStart()
	}
	pc.gcRoutineStart()
	return pc, nil
}

func (me *ValueBehaviour) check() {
	if me.ValueBeType != None {
		if me.callOnExpire() {
			if me.OnExpireFunc == nil {
				if me.OnExpireMethodName == "" {
					me.OnExpireMethodName = onExpire
				}
			}
		}
		if me.callOnLoad() {
			if me.OnLoadFunc == nil {
				if me.OnLoadMethodName == "" {
					me.OnLoadMethodName = onLoad
				}
			}
		}
	}
}

func (me *ValueBehaviour) callOnLoad() bool {
	return me.ValueBeType == CallOnLoadAndOnExpire
}

func (me *ValueBehaviour) callOnLoadWithReflection() bool {
	return me.OnLoadMethodName != ""
}

func (me *ValueBehaviour) callOnExpireWithReflection() bool {
	return me.OnExpireMethodName != ""
}

func (me *ValueBehaviour) callOnExpire() bool {
	return me.ValueBeType == CallOnExpire || me.ValueBeType == CallOnLoadAndOnExpire
}

func (me *UCache) diskSyncStart() {
	if me.speedMode {
		me.diskSyncStopped = false
		me.diskSyncStopSignal = make(chan bool)
		me.diskSyncNotify = make(chan bool, 200)
		persistAllMemCacheRunning := false
		persistAllSignal := make(chan bool, 1)
		persistAllSignal <- true
		go func() {
			defer func() {
				if persistAllMemCacheRunning {
					<-persistAllSignal
				}
				close(me.diskSyncStopSignal)
			}()
			for {
				select {
				case <-me.diskSyncNotify:
					if me.diskCache && !me.diskSyncStopped && !persistAllMemCacheRunning {
						persistAllMemCacheRunning = true
						<-persistAllSignal
						go func() {
							me.persistAllMemCache()
							persistAllSignal <- true
							persistAllMemCacheRunning = false
						}()
					}
				case <-me.diskSyncStopSignal:
					return
				}
			}

		}()
	}
}

func (me *UCache) diskSyncStop() {
	if me.speedMode && !me.diskSyncStopped {
		me.diskSyncStopped = true
		me.diskSyncStopSignal <- true
		<-me.diskSyncStopSignal
	}
}

func (me *UCache) gcRoutineStart() {
	if me.noExpiry || !me.gcStopped {
		return
	}

	me.gcStoppedLock.Lock()
	me.gcStopped = false
	me.gcStoppedLock.Unlock()
	if me.gcStop != nil {
		_, open := <-me.gcStop
		if open {
			panic(fmt.Errorf("chan still open"))
		}
	}

	me.gcStop = make(chan bool, 1)
	go func() {
		lastExpiry := me.expiry
		minExpiry := lastExpiry
		ticker := time.NewTicker(lastExpiry)
		defer func() {
			me.gcStoppedLock.Lock()
			me.gcStopped = true
			me.gcStoppedLock.Unlock()
			close(me.gcStop)
		}()
		updateTimer := func(t time.Duration) {
			if lastExpiry != t {
				ticker.Stop()
				ticker = time.NewTicker(t)
				lastExpiry = t
			}
		}
		minExpiry = me.gc()
		if minExpiry > 1 {
			updateTimer(minExpiry)
		} else {
			updateTimer(me.expiry)
		}
		for {
			select {
			case <-ticker.C:
				minExpiry = me.gc()
				if minExpiry > 1 {
					updateTimer(minExpiry)
				} else {
					updateTimer(me.expiry)
				}
			case <-me.gcStop:
				ticker.Stop()
				return
			}
		}
	}()
}
func (me *UCache) gcRoutineStop() {
	if me.noExpiry {
		return
	}
	me.gcStoppedLock.Lock()
	stopped := me.gcStopped
	me.gcStoppedLock.Unlock()
	if !stopped {
		me.gcStop <- true
		<-me.gcStop
		me.nextToBeCleanedUp = make([][]byte, 0)
	}
}

func (me *UCache) DiskPath() string {
	return me.diskStorePath
}

func (me *UCache) DiskFilename() string {
	return me.diskStoreName
}

func (me *UCache) DiskDir() string {
	return me.diskStoreDir
}

func (me *UCache) openBoltDB() (err error) {
	me.diskStore, err = bolt.Open(me.diskStorePath, 0600, &bolt.Options{NoGrowSync: true, Timeout: 6 * time.Second})
	return
}

func (me *UCache) checkPath() {
	me.diskStorePath = strings.TrimSpace(me.diskStorePath)
	me.diskStoreName = "."
	if !strings.HasSuffix(me.diskStorePath, string(os.PathSeparator)) {
		me.diskStoreName = filepath.Base(me.diskStorePath)
	}
	if me.diskStoreName == "." {
		me.diskStoreName = "cache"
	}
	me.diskStoreDir = filepath.Dir(me.diskStorePath)
	me.diskStorePath = filepath.Join(me.diskStoreDir, me.diskStoreName)
}

func (me *UCache) ensureDir() error {
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

func (me *UCache) pathExists() bool {
	_, err := os.Stat(me.diskStorePath)
	return !os.IsNotExist(err)
}

func (me *UCache) IterateMemStorage(clb func(key string, val interface{})) {
	me.memLock.RLock()
	for k, v := range me.memStore {
		clb(k, v.val)
	}
	me.memLock.RUnlock()
}

func (me *UCache) Get(key string, ref interface{}) error {
	return me.GetAndExtendExpiry(key, ref, me.defaultExtendExpiry)
}

func (me *UCache) GetAndExtendExpiry(key string, ref interface{}, extendExpiry bool) error {
	if me.memCache {
		me.memLock.RLock()
		memvh := me.memStore[key]
		me.memLock.RUnlock()

		if memvh != nil {
			if me.diskCache && memvh.loadingFromDisk && me.storeType != DiskCacheSpeedMode {
				//wait in case it is already loading from disk
				memvh.m.Lock()
				if memvh.loadingFromDiskErr == ErrNotFound {
					memvh.m.Unlock()
					return memvh.loadingFromDiskErr
				}
				memvh.m.Unlock()
			}

			//update last touch
			if extendExpiry {
				memvh.m.Lock()
				memvh.access = time.Now()
				memvh.persisted = false
				memvh.m.Unlock()

				if me.diskCache && !me.speedMode {
					meta := diskMeta{A: memvh.access, D: memvh.duration}
					metaBts, err := json.Marshal(meta)
					if err != nil {
						return err
					}
					err = me.diskPut(&key, metaBts, nil)
					if err != nil {
						return err
					}
				}
			}

			return me.assignValToRef(memvh, key, ref, extendExpiry)
		}
	}

	if me.diskCache {
		// load from disk
		var err error
		var val *value
		if me.memCache && me.storeType != DiskCacheSpeedMode {
			iAmLoadingFromDisk := false
			//keep others waiting for this key while it is already loading from disk
			me.memLock.Lock()
			if v, ok := me.memStore[key]; ok {
				val = v
				me.memLock.Unlock()
				val.m.Lock()
			} else {
				iAmLoadingFromDisk = true
				val = &value{
					loadingFromDisk: true,
					persisted:       true,
					access:          time.Now(),
					duration:        me.expiry,
				}
				me.memStore[key] = val
				//overlap to ensure we as the loader are keeping the monitor
				//no deadlock risk as we just created the value
				val.m.Lock()
				me.memLock.Unlock()
			}

			defer func() {
				//end of the call, resetting flag and unlocking
				//delete the value in case loading from disk returned an error
				if iAmLoadingFromDisk {
					val.loadingFromDisk = false
					if err != nil {
						val.loadingFromDiskErr = err
						val.m.Unlock()
						me.memLock.Lock()
						delete(me.memStore, key)
						me.memLock.Unlock()
						return
					}
				}
				val.m.Unlock()

			}()
			if !iAmLoadingFromDisk {
				if val.loadingFromDiskErr == ErrNotFound {
					//value was very recently searched on the disk but was not found
					return val.loadingFromDiskErr
				} else if val.loadingFromDiskErr == nil {
					//value was very recently loaded from the disk and already set on memory cache
					return me.assignValToRef(val, key, ref, extendExpiry)
				}
			}
		}

		meta := diskMeta{}
		//--------------- loading from disk ---------------------
		err = me.diskStore.View(func(tx *bolt.Tx) error {
			b1 := tx.Bucket(bucketMeta)
			b2 := tx.Bucket(bucketData)
			if b1 != nil && b2 != nil {
				metaBts := b1.Get([]byte(key))
				if len(metaBts) == 0 {
					return ErrNotFound
				}
				err = json.Unmarshal(metaBts, &meta)
				if err != nil {
					return err
				}
				err = json.Unmarshal(b2.Get([]byte(key)), ref)
				if err != nil {
					return err
				}
				return nil
			}
			return ErrNotFound
		})
		if err != nil {
			return err
		}
		//-------------------------------

		//update last touch
		if extendExpiry {
			if val != nil {
				meta.A = val.access
			} else {
				meta.A = time.Now()
			}
			var metaBts []byte
			metaBts, err = json.Marshal(meta)
			if err != nil {
				return err
			}
			err = me.diskPut(&key, metaBts, nil)
			if err != nil {
				return err
			}
		}

		if me.memCache && me.storeType != DiskCacheSpeedMode {
			//vv := reflect.Indirect(reflect.ValueOf(ref))
			vv := reflect.ValueOf(ref).Elem()
			if vv.CanInterface() {
				//set loaded value from disk to memCache
				if !extendExpiry {
					val.access = meta.A
				}
				val.duration = meta.D
				val.val = vv.Interface()
				//after setting the value to memCache, we check if we have to call OnLoad on it
				if me.valueBehaviour.callOnLoad() {
					if me.valueBehaviour.callOnLoadWithReflection() {
						CallMethodIfExists(val.val, me.valueBehaviour.OnLoadMethodName)
					}
					if me.valueBehaviour.OnLoadFunc != nil {
						me.valueBehaviour.OnLoadFunc(key, val.val)
					}
				}
			}
			return nil
		}
		//as we are not allowed to touch the memCache here
		//we check only if we need to call OnLoad
		if me.valueBehaviour.callOnLoad() {
			vv := reflect.Indirect(reflect.ValueOf(ref))
			if vv.CanInterface() {
				if me.valueBehaviour.callOnLoadWithReflection() {
					CallMethodIfExists(vv.Interface(), me.valueBehaviour.OnLoadMethodName)
				}
				if me.valueBehaviour.OnLoadFunc != nil {
					me.valueBehaviour.OnLoadFunc(key, vv.Interface())
				}
			}
		}
		return nil
	}
	return ErrNotFound
}

func (me *UCache) assignValToRef(memvh *value, key string, ref interface{}, extendExpiry bool) error {
	v := reflect.ValueOf(ref)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return os.ErrInvalid
	}

	memvh.m.RLock()
	defer memvh.m.RUnlock()
	mv := reflect.ValueOf(memvh.val)
	if !mv.IsValid() {
		return ErrNotFound
	}
	i := 0
	for v.Kind() != reflect.Struct && v.Kind() != reflect.Invalid && (!v.CanSet() || v.Type() != mv.Type()) {
		v = v.Elem()
		if i > 3 {
			break
		}
		i++
	}
	if !v.CanSet() || v.Kind() != mv.Kind() {
		return os.ErrInvalid
	}
	v.Set(mv)
	return nil
}

func (me *UCache) diskPut(key *string, meta []byte, data []byte) error {
	return me.diskStore.Update(func(tx *bolt.Tx) error {
		if len(meta) > 0 {
			b, err := tx.CreateBucketIfNotExists(bucketMeta)
			if err != nil {
				return err
			}
			err = b.Put([]byte(*key), meta)
			if err != nil {
				return err
			}
		}
		if len(data) > 0 {
			b, err := tx.CreateBucketIfNotExists(bucketData)
			if err != nil {
				return err
			}
			err = b.Put([]byte(*key), data)
			return err
		}
		return nil
	})
}

func (me *UCache) Remove(key string) error {
	if me.memCache {
		me.memLock.Lock()
		if v, ok := me.memStore[key]; ok {
			me.callOnExpireOnValInner(&key, v)
		}
		delete(me.memStore, key)
		me.memLock.Unlock()
	}
	if me.diskCache {
		return me.diskStore.Update(func(tx *bolt.Tx) error {
			var err error
			b := tx.Bucket(bucketMeta)
			if b != nil {
				err = b.Delete([]byte(key))
				if err != nil {
					return err
				}
			}
			b = tx.Bucket(bucketData)
			if b != nil {
				err = b.Delete([]byte(key))
			}
			return err
		})
	}
	return nil
}

func (me *UCache) Put(key string, val interface{}) error {
	return me.PutWithOtherExpiry(key, val, me.expiry)
}

func (me *UCache) PutWithOtherExpiry(key string, val interface{}, differentExpiry time.Duration) error {
	n := time.Now()
	me.lastPut = n
	if me.memCache {
		me.memLock.Lock()
		v := me.memStore[key]
		if v == nil {
			v = &value{
				persisted: false,
				duration:  differentExpiry,
				access:    n,
				val:       val,
			}
			me.memStore[key] = v
			me.memLock.Unlock()
		} else {
			me.memLock.Unlock()
			v.m.Lock()
			v.access = n
			v.duration = differentExpiry
			v.persisted = false
			v.val = val
			v.m.Unlock()
		}
		if me.speedMode {
			me.diskSyncNotify <- true
		}
	}

	if me.diskCache && !me.speedMode {
		meta := diskMeta{A: n, D: differentExpiry}
		metaBts, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		dataBts, err := json.Marshal(val)
		if err != nil {
			return err
		}
		return me.diskPut(&key, metaBts, dataBts)
	}
	return nil
}

func (me *UCache) UpdatedValueRef(key string) error {
	if !me.memCache {
		return ErrNotSupportedWithoutMemCache
	}
	n := time.Now()
	me.lastPut = n

	me.memLock.RLock()
	v := me.memStore[key]
	me.memLock.RUnlock()

	if v == nil {
		return ErrNotFound
	}

	v.m.Lock()
	v.access = n
	v.persisted = false
	v.m.Unlock()

	if me.diskCache && !me.speedMode {
		meta := diskMeta{A: n, D: v.duration}
		metaBts, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		dataBts, err := json.Marshal(v.val)
		if err != nil {
			return err
		}
		return me.diskPut(&key, metaBts, dataBts)
	}
	return nil
}

func (me *UCache) gc() time.Duration {
	if me.diskCache {
		if me.diskCacheSize() == 0 {
			return -1
		}
	} else {
		me.memLock.RLock()
		memSize := len(me.memStore)
		me.memLock.RUnlock()
		if memSize == 0 {
			return -1
		}
	}
	if me.verbose {
		log.Println("running gc")
	}
	n := time.Now()
	me.lastGCRun = n
	minExpiry := n.Add(me.expiry)
	cacheCount := 0
	expiredKeys := make([]string, 0)
	if me.diskCache {
		err := me.diskStore.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketMeta)
			if b != nil {
				meta := diskMeta{}
				percent95ofExpiry := me.expiry / time.Duration(100) * time.Duration(95)
				expiry95Percent := n.Add(me.expiry).Add(percent95ofExpiry * time.Duration(-1))
				if len(me.nextToBeCleanedUp) > 0 {
					for _, k := range me.nextToBeCleanedUp {
						v := b.Get(k)
						err := json.Unmarshal(v, &meta)
						if err == nil {
							exp := meta.expiry()
							if time.Now().After(exp) {
								expiredKeys = append(expiredKeys, string(k))
							} else if exp.Before(minExpiry) {
								minExpiry = exp
							}
						}
					}
					me.nextToBeCleanedUp = make([][]byte, 0)
				} else {
					const arrayCap = 200000
					startedLongRunningTaskTime := time.Now()
					nextCount := 0
					expCount := 0
					runningAlreadyForToLong := func(l int) bool {
						//ensure we are not filling up the memory for cleaning up to much
						if l > arrayCap {
							if me.verbose {
								log.Println("gc collection on cap.. breaking up: collected entries", l)
							}
							return true
						}
						if time.Now().After(startedLongRunningTaskTime.Add(time.Millisecond * 1500)) {
							if me.verbose {
								log.Println("gc took longer than 1.5 sec.. breaking up: collected entries", l)
							}
							return true
						}
						return false
					}
					c := b.Cursor()
					for k, v := c.First(); k != nil; k, v = c.Next() {
						err := json.Unmarshal(v, &meta)
						kcopy := make([]byte, len(k))
						copy(kcopy, k)
						if err == nil {
							exp := meta.expiry()
							if time.Now().After(exp) {
								expiredKeys = append(expiredKeys, string(kcopy))
								expCount++
								if runningAlreadyForToLong(expCount) {
									break
								}
							} else if exp.Before(expiry95Percent) {
								if exp.Before(minExpiry) {
									minExpiry = exp
								}
								me.nextToBeCleanedUp = append(me.nextToBeCleanedUp, kcopy)
								nextCount++
								if runningAlreadyForToLong(nextCount) {
									break
								}
							}
						}
					}
				}
			}
			return nil
		})
		diskErrorPrint(err)
		if len(expiredKeys) > 0 {
			meta := diskMeta{}
			err = me.diskStore.Update(func(tx *bolt.Tx) error {
				b1 := tx.Bucket(bucketMeta)
				b2 := tx.Bucket(bucketData)
				if b1 != nil && b2 != nil {
					var er error
					size := len(expiredKeys)
					for i := 0; i < size; i++ {
						er = nil
						k := []byte(expiredKeys[i])
						if !me.metaFromMemCache(&expiredKeys[i], &meta) {
							er = json.Unmarshal(b1.Get(k), &meta)
						}
						if er == nil {
							exp := meta.expiry()
							if time.Now().After(exp) {
								er = b1.Delete(k)
								diskErrorPrint(er)
								er = b2.Delete(k)
								diskErrorPrint(er)
								if me.memCache {
									if me.valueBehaviour.callOnExpire() {
										me.callOnExpireOnVal(&expiredKeys[i])
									}
								} else if me.valueBehaviour.OnExpireFunc != nil {
									me.valueBehaviour.OnExpireFunc(expiredKeys[i], nil)
								}
								continue
							}
						}
						expiredKeys = append(expiredKeys[:i], expiredKeys[i+1:]...)
						size--
						i--
					}
				}
				return err
			})
			diskErrorPrint(err)
			cacheCount = me.diskTryDeleteDBIfEmpty()
			if me.memCache && me.storeType != DiskCacheSpeedMode {
				me.memLock.Lock()
				for _, key := range expiredKeys {
					delete(me.memStore, key)
				}
				me.memLock.Unlock()
			}
		}

	} else if me.memCache {
		me.memLock.Lock()
		for key, val := range me.memStore {
			val = me.memStore[key]
			exp := val.expiry()
			if n.After(exp) {
				expiredKeys = append(expiredKeys, key)
			} else if exp.Before(minExpiry) {
				minExpiry = exp
			}
		}
		if me.valueBehaviour.ValueBeType == None {
			for _, key := range expiredKeys {
				delete(me.memStore, key)
			}
		}
		cacheCount = len(me.memStore)
		me.memLock.Unlock()
		if me.valueBehaviour.callOnExpire() && (me.valueBehaviour.callOnExpireWithReflection() || me.valueBehaviour.OnExpireFunc != nil) {
			for _, key := range expiredKeys {
				me.callOnExpireOnVal(&key)
			}
		}
		if me.valueBehaviour.ValueBeType != None {
			me.memLock.Lock()
			for _, key := range expiredKeys {
				delete(me.memStore, key)
			}
			me.memLock.Unlock()
		}
	}
	if me.verbose {
		log.Println("gc done current cache count", cacheCount, "expired count", len(expiredKeys))
	}
	return minExpiry.Sub(n)
}

func (me *UCache) metaFromMemCache(key *string, meta *diskMeta) bool {
	if me.memCache {
		me.memLock.RLock()
		v := me.memStore[*key]
		me.memLock.RUnlock()
		if v != nil {
			meta.A = v.access
			meta.D = v.duration
			return true
		}
	}
	return false
}

func (me *UCache) callOnExpireOnVal(key *string) {
	me.memLock.RLock()
	v := me.memStore[*key]
	me.memLock.RUnlock()

	me.callOnExpireOnValInner(key, v)
}

func (me *UCache) callOnExpireOnValInner(key *string, v *value) {
	if me.valueBehaviour.callOnExpireWithReflection() && v != nil {
		CallMethodIfExists(v.val, me.valueBehaviour.OnExpireMethodName)
	}
	if me.valueBehaviour.OnExpireFunc != nil {
		if v != nil {
			me.valueBehaviour.OnExpireFunc(*key, v.val)
		} else {
			me.valueBehaviour.OnExpireFunc(*key, v)
		}
	}
}

func (me *UCache) diskTryDeleteDBIfEmpty() (cacheCount int) {
	cacheCount = me.diskCacheSize()
	if me.verbose {
		log.Println("cacheCount", cacheCount)
	}

	if cacheCount == 0 {
		if me.verbose {
			log.Println("try delete db")
		}
		//remove db file as boltdb is not freeing space after calling delete
		me.memCache = true
		me.diskCache = false
		me.diskSyncLock.Lock()
		defer me.diskSyncLock.Unlock()
		cacheCount = me.diskCacheSize()
		if cacheCount == 0 {
			me.diskStore.Close()
			os.Remove(me.diskStorePath)
			err := me.openBoltDB()
			//TODO if err != nil impl recover like "keep memCache true until me.openBoltDB() is not returning errors and write memCache to diskCache" or panic
			me.diskCache = true
			if me.verbose {
				log.Println("deleted db")
			}
			diskErrorPrint(err)
			if !me.speedMode {
				me.persistAllMemCacheAndRevertStoreFlags()
			}
		} else {
			me.diskCache = true
			if !me.speedMode {
				me.persistAllMemCacheAndRevertStoreFlags()
			}
		}
	}
	return
}

func (me *UCache) persistAllMemCache() {
	if me.verbose {
		log.Println("persist all mem cache", len(me.memStore))
	}
	me.diskSyncLock.Lock()
	defer func() {
		me.diskSyncLock.Unlock()
		if me.verbose {
			log.Println("persist all mem cache done")
		}
	}()
	if !me.diskCache || me.diskSyncStopped {
		return
	}

	for {
		if !me.diskCache || me.diskSyncStopped {
			return
		}
		done, err := me.persistBatch()
		if err != nil {
			diskErrorPrint(err)
			return
		}
		if !me.diskCache || me.diskSyncStopped {
			return
		}
		if me.storeType == DiskCacheSpeedMode {
			me.freePersistedMemCache()
		}
		if !me.diskCache || me.diskSyncStopped {
			return
		}
		if done {
			runtime.GC()
			break
		}
	}
}

func (me *UCache) persistBatch() (bool, error) {
	// Start the transaction.
	tx, err := me.diskStore.Begin(true)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	// Setup the users bucket.
	b1, err := tx.CreateBucketIfNotExists(bucketMeta)
	if err != nil {
		return false, err
	}
	b2, err := tx.CreateBucketIfNotExists(bucketData)
	if err != nil {
		return false, err
	}
	maxMBPerTx := 2
	byteSizePerTx := 0
	entryCount := 0
	allEntriesPersisted := true
	minExpiry := time.Now()
	chunkSize := len(me.memStore)
	if chunkSize > 200 {
		chunkSize = 200
	}
	for {
		metas := me.iterateNotPersistedMemCache(chunkSize)
		if len(metas) == 0 {
			break
		}
		for key, val := range metas {
			if !me.diskCache || me.diskSyncStopped {
				tx.Commit()
				return false, nil
			}
			meta := diskMeta{A: val.access, D: val.duration}
			if meta.before(minExpiry) {
				minExpiry = meta.expiry()
			}
			metaBts, err := json.Marshal(meta)
			if err == nil {
				dataBts, err := json.Marshal(val.val)
				if err == nil {
					k := []byte(key)
					byteSizePerTx += len(k)
					byteSizePerTx += len(k)
					byteSizePerTx += len(metaBts)
					byteSizePerTx += len(dataBts)
					err = b1.Put(k, metaBts)
					if err != nil {
						tx.Commit()
						return false, err
					}
					err = b2.Put(k, dataBts)
					if err != nil {
						tx.Commit()
						return false, err
					}
					val.m.Lock()
					val.persisted = true
					val.m.Unlock()
					entryCount++
					if byteSizePerTx/1024/1024 >= maxMBPerTx {
						allEntriesPersisted = false
						break
					}
				} else {
					log.Println("marshal error: ", err)
				}
			} else {
				log.Println("marshal error: ", err)
			}
		}
		if !allEntriesPersisted {
			break
		}
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return false, err
	}
	if me.verbose {
		log.Println("committed entries", entryCount)
	}
	return allEntriesPersisted, nil
}

func (me *UCache) iterateNotPersistedMemCache(size int) map[string]*value {
	metas := make(map[string]*value, size)
	i := 0
	n := time.Now()
	me.memLock.RLock()
	for key, val := range me.memStore {
		if !val.persisted {
			if i >= size {
				break
			}
			//check if already expired
			if n.After(val.expiry()) {
				val.m.Lock()
				val.persisted = true
				val.m.Unlock()
				continue
			}
			metas[key] = val
			i++
		}
	}
	me.memLock.RUnlock()
	if i > 0 {
		return metas
	}
	return nil
}

func (me *value) String() string {
	return fmt.Sprintf(`{"loadingFromDisk":%v,"persisted":%v, "access":"%v", "duration":"%v"}`, me.loadingFromDisk, me.persisted, me.access, me.duration)
}

func (me *UCache) freePersistedMemCache() {
	count := 0
	me.memLock.Lock()
	for key, val := range me.memStore {
		if val.persisted {
			count++
			delete(me.memStore, key)
		}
	}
	me.memLock.Unlock()
	if me.verbose {
		log.Println("freed", count)
	}
}

func (me *UCache) persistAllMemCacheAndRevertStoreFlags() {
	me.persistAllMemCache()
	me.revertStoreFlags()
}

func (me *UCache) revertStoreFlags() {
	if me.storeType == DiskCache {
		me.memLock.Lock()
		me.memCache = false
		me.memStore = make(map[string]*value)
		me.memLock.Unlock()
	}
}

func diskErrorPrint(err error) {
	if err != nil {
		log.Println("disk error when cleaning up cache: ", err)
	}
}

func toByteSlice(a uint64) []byte {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(a))
	return bs
}

//Size correctness is not guaranteed
func (me *UCache) Size() int {
	size := 0
	if me.memCache {
		me.memLock.RLock()
		size = len(me.memStore)
		me.memLock.RUnlock()
	}
	if size == 0 && me.diskCache {
		return me.diskCacheSize()
	} else if me.diskCache && !me.speedMode {
		return me.diskCacheSize()
	}
	return size
}

func (me *UCache) diskCacheSize() int {
	size := 0
	err := me.diskStore.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMeta)
		if b != nil {
			size = b.Stats().KeyN
		}
		return nil
	})
	if err == bolt.ErrDatabaseNotOpen {
		return size
	}
	diskErrorPrint(err)
	return size
}

//Clean removes disk and mem cache and makes everything ready to start from scratch.
//OnExpiry method is called on all entries in mem cache if configured even though they might have some time left until expiry
func (me *UCache) Clean() error {
	me.gcRoutineStop()
	if me.diskCache {
		me.diskCache = false
		me.diskSyncStop()
		me.diskStore.Close()
		os.Remove(me.diskStorePath)
		err := me.openBoltDB()
		if err != nil {
			return err
		}
		me.diskCache = true
	}
	me.expireAllMem()
	me.memLock.Lock()
	me.memStore = make(map[string]*value)
	me.memLock.Unlock()
	if me.diskCache && me.speedMode {
		me.diskSyncStart()
	}
	me.gcRoutineStart()
	return nil
}

func (me diskMeta) expired() bool {
	return time.Now().After(me.A.Add(me.D))
}
func (me diskMeta) expiry() time.Time {
	return me.A.Add(me.D)
}
func (me diskMeta) before(t time.Time) bool {
	return me.A.Add(me.D).Before(t)
}

func (me *value) expiry() time.Time {
	return me.access.Add(me.duration)
}

func (me StoreType) String() string {
	switch me {
	case DiskCacheSpeedMode:
		return "DiskCacheSpeedMode"
	case MemCache:
		return "MemCache"
	case DiskCache:
		return "DiskCache"
	case MemAndDiskCache:
		return "MemAndDiskCache"
	case MemAndDiskCacheSpeedMode:
		return "MemAndDiskCacheSpeedMode"
	}
	return ""
}

var zero = reflect.Value{}
var zeroArgs = make([]reflect.Value, 0)

func CallMethodIfExists(obj interface{}, mname string) {
	vr := reflect.ValueOf(obj)
	//---- interface{} != nil not in any case supported in go
	if strings.HasPrefix(vr.Type().String(), "*") {
		if vr.IsNil() {
			return
		}
	}
	//-----
	if obj != nil && mname != "" {
		method := vr.MethodByName(mname)
		if method != zero && method.Type().NumIn() == 0 {
			method.Call(zeroArgs)
		}
	}
}

func (me *UCache) expireAllMem() {
	if me.memCache {
		if me.valueBehaviour.callOnExpire() && (me.valueBehaviour.callOnExpireWithReflection() || me.valueBehaviour.OnExpireFunc != nil) {
			me.memLock.Lock()
			for key, val := range me.memStore {
				me.callOnExpireOnValInner(&key, val)
			}
			me.memLock.Unlock()
		}
	}
}

func (me *UCache) Close() {
	if me.verbose {
		log.Println("Close")
	}
	me.gcRoutineStop()
	me.diskSyncStop()
	if me.diskCache {
		if me.speedMode {
			me.diskSyncStopped = false
			me.persistAllMemCache()
		}
		me.diskStore.Close()
	}
}
