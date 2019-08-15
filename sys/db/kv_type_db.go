package db

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type (
	KVStore struct {
		store       KVStoreIF              `json:"-"`
		data        map[string]interface{} `json:"-"`
		dataLock    sync.RWMutex           `json:"-"`
		memData     map[string]interface{} `json:"-"`
		memDataLock sync.RWMutex           `json:"-"`
	}

	KVStoreIF interface {
		Create(path string) error
		Put(key *string, val []byte) error
		Get(key *string) ([]byte, error)
		Delete(key *string) error
		All() (keys []string, err error)
		Close(delete bool) error
	}

	ErrInvalidType struct {
		err      error
		TypeExp  reflect.Type
		TypeProv reflect.Type
	}
)

func NewKVStore(storeIF KVStoreIF, path string) (*KVStore, error) {
	if storeIF != nil {
		err := storeIF.Create(path)
		if err != nil {
			return nil, err
		}
	}
	return &KVStore{
		store:   storeIF,
		data:    map[string]interface{}{},
		memData: map[string]interface{}{}}, nil
}

func (me *KVStore) Get(key string, ref interface{}) (err error) {
	v := reflect.ValueOf(ref)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return &ErrInvalidType{
			TypeProv: v.Type(),
		}
	}
	var ok bool
	var d interface{}

	// load from mem only
	me.memDataLock.RLock()
	d, ok = me.memData[key]
	me.memDataLock.RUnlock()

	if ok {
		i := 0
		for v.Kind() != reflect.Struct && v.Kind() != reflect.Invalid && (!v.CanSet() || v.Type() != reflect.TypeOf(d)) {
			v = v.Elem()
			if i > 3 {
				break
			}
			i++
		}
		if !v.CanSet() || v.Kind() != reflect.TypeOf(d).Kind() {
			return &ErrInvalidType{
				TypeExp:  reflect.TypeOf(d),
				TypeProv: reflect.TypeOf(v),
			}
		}
		v.Set(reflect.ValueOf(d))
		return
	}

	// load from mem
	me.dataLock.RLock()
	d, ok = me.data[key]
	me.dataLock.RUnlock()
	if ok {
		i := 0
		for v.Kind() != reflect.Struct && v.Kind() != reflect.Invalid && (!v.CanSet() || v.Type() != reflect.TypeOf(d)) {
			v = v.Elem()
			if i > 3 {
				break
			}
			i++
		}
		if !v.CanSet() || v.Kind() != reflect.TypeOf(d).Kind() {
			return &ErrInvalidType{TypeExp: reflect.TypeOf(d), TypeProv: reflect.TypeOf(v)}
		}
		v.Set(reflect.ValueOf(d))
		return
	}

	// load from disk
	if me.store != nil {
		var bts []byte
		if err != nil {
			return err
		}
		bts, err = me.store.Get(&key)
		if err != nil {
			return err
		}
		err = json.Unmarshal(bts, ref)
		vv := reflect.Indirect(reflect.ValueOf(ref))
		if vv.CanInterface() {
			//set to mem
			me.dataLock.Lock()
			me.data[key] = vv.Interface()
			me.dataLock.Unlock()
		}
	}
	return
}

func (me *KVStore) Put(key string, ref interface{}) (err error) {
	me.dataLock.Lock()
	defer me.dataLock.Unlock()
	me.data[key] = ref
	if me.store != nil {
		var bts []byte
		bts, err = json.Marshal(ref)
		if err != nil {
			return err
		}
		return me.store.Put(&key, bts)
	}
	return nil
}

func (me *KVStore) PutMemOnly(key string, ref interface{}) (err error) {
	me.memDataLock.Lock()
	defer me.memDataLock.Unlock()
	me.memData[key] = ref
	return nil
}

func (me *KVStore) Delete(key string) (err error) {
	me.dataLock.Lock()
	defer me.dataLock.Unlock()
	if val, ok := me.data[key]; ok {
		me.callMethodIfExists(val, "Close")
	}
	delete(me.data, key)
	if me.store != nil {
		return me.store.Delete(&key)
	}
	return nil
}

func (me *KVStore) All() (keys []string, err error) {
	if me.store != nil {
		return me.store.All()
	}
	return
}

func (me *KVStore) Close(delete bool) (err error) {
	if me.store != nil {
		if delete {
			me.memDataLock.Lock()
			defer me.memDataLock.Unlock()
			for _, val := range me.memData {
				me.callMethodIfExists(val, "Close")
			}
		}
		me.dataLock.Lock()
		defer me.dataLock.Unlock()
		for name, val := range me.data {
			if delete {
				me.callMethodIfExists(val, "Close")
			} else {
				var bts []byte
				bts, err = json.Marshal(val)
				if err != nil {
					return err
				}
				me.store.Put(&name, bts)
			}
		}
		return me.store.Close(delete)
	}
	return nil
}

var zero = reflect.Value{}

func (me *KVStore) callMethodIfExists(obj interface{}, mname string) {
	vr := reflect.ValueOf(obj)
	//---- interface{} != nil not in any case supported in go yet
	if strings.HasPrefix(reflect.TypeOf(obj).String(), "*") {
		if vr.IsNil() {
			return
		}
	}
	//-----
	if obj != nil && mname != "" {
		method := vr.MethodByName(strings.Title(mname))
		if method != zero && method.Type().NumIn() == 0 {
			method.Call([]reflect.Value{})
		}
	}
}

func (v ErrInvalidType) Error() string {
	if v.err == nil {
		if v.TypeExp == nil {
			v.err = fmt.Errorf("expected: *interface{}, provided: %v", v.TypeProv)
		} else {
			v.err = fmt.Errorf("expected: %v, provided: %v", v.TypeExp, v.TypeProv)
		}
	}
	return v.err.Error()
}
