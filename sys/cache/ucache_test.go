package cache

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var config = UCacheConfig{
	DiskStorePath: "testCache",
	ExtendExpiry:  false,
	DefaultExpiry: 30 * time.Millisecond,
	Verbose:       true,
}

func TestDiskCacheSpeedMode(t *testing.T) {
	config.StoreType = DiskCacheSpeedMode
	c := ioTest(&config, t)
	os.Remove(c.DiskPath())
}

func TestMemCache(t *testing.T) {
	config.StoreType = MemCache
	ioTest(&config, t)
}

func TestDiskCache(t *testing.T) {
	config.StoreType = DiskCache
	c := ioTest(&config, t)
	os.Remove(c.DiskPath())
}

func TestMemAndDiskCache(t *testing.T) {
	config.StoreType = MemAndDiskCache
	c := ioTest(&config, t)
	os.Remove(c.DiskPath())
}

func TestMemAndDiskCacheSpeedMode(t *testing.T) {
	config.StoreType = MemAndDiskCacheSpeedMode
	c := ioTest(&config, t)
	os.Remove(c.DiskPath())
}

func TestDiskCacheSpeedModeCallOnLoadAndOnExpire(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnExpireMethodName: onExpire}
	config.StoreType = DiskCacheSpeedMode
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestDiskCacheSpeedModeCallOnLoadAndClose(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnExpireMethodName: "Close"}
	config.StoreType = DiskCacheSpeedMode
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestDiskCacheSpeedModeCallOnExpire(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnExpire, OnExpireMethodName: onExpire}
	config.StoreType = DiskCacheSpeedMode
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestDiskCacheSpeedModeCallClose(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnExpire, OnExpireMethodName: "Close"}
	config.StoreType = DiskCacheSpeedMode
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestMemCacheCallOnLoadAndOnExpire(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnExpireMethodName: onExpire}
	config.StoreType = MemCache
	ioLifeCycleTest(&config, t, nil)
}

func TestMemCacheCallOnLoadAndClose(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnExpireMethodName: "Close"}
	config.StoreType = MemCache
	ioLifeCycleTest(&config, t, nil)
}

func TestMemCacheCallOnExpire(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnExpire, OnExpireMethodName: onExpire}
	config.StoreType = MemCache
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestMemCacheCallClose(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnExpire, OnExpireMethodName: "Close"}
	config.StoreType = MemCache
	ioLifeCycleTest(&config, t, nil)
}

func TestDiskCacheCallOnLoadAndOnExpire(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnExpireMethodName: onExpire}
	config.StoreType = DiskCache
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestDiskCacheCallOnLoadAndClose(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnExpireMethodName: "Close"}
	config.StoreType = DiskCache
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestDiskCacheCallOnExpire(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnExpire, OnExpireMethodName: onExpire}
	config.StoreType = DiskCache
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestDiskCacheCallClose(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnExpire, OnExpireMethodName: "Close"}
	config.StoreType = DiskCache
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestMemAndDiskCacheCallOnLoadAndOnExpire(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnExpireMethodName: onExpire}
	config.StoreType = MemAndDiskCache
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestMemAndDiskCacheCallOnLoadAndClose(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnExpireMethodName: "Close"}
	config.StoreType = MemAndDiskCache
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestMemAndDiskCacheCallOnExpire(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnExpire, OnExpireMethodName: onExpire}
	config.StoreType = MemAndDiskCache
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestMemAndDiskCacheCallClose(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnExpire, OnExpireMethodName: "Close"}
	config.StoreType = MemAndDiskCache
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestMemAndDiskCacheSpeedModeCallOnLoadAndOnExpire(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnExpireMethodName: onExpire}
	config.StoreType = MemAndDiskCacheSpeedMode
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestMemAndDiskCacheSpeedModeCallOnLoadAndClose(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnExpireMethodName: "Close"}
	config.StoreType = MemAndDiskCacheSpeedMode
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestMemAndDiskCacheSpeedModeCallOnExpire(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnExpire, OnExpireMethodName: onExpire}
	config.StoreType = MemAndDiskCacheSpeedMode
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestMemAndDiskCacheSpeedModeCallClose(t *testing.T) {
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnExpire, OnExpireMethodName: "Close"}
	config.StoreType = MemAndDiskCacheSpeedMode
	c := ioLifeCycleTest(&config, t, nil)
	os.Remove(c.DiskPath())
}

func TestMemAndDiskCacheSpeedModeCallOnLoadAndOnExpireFuncs(t *testing.T) {
	w := &wrapper{}
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnLoadFunc: w.OnLoad, OnExpireFunc: w.OnExpire}
	config.StoreType = MemAndDiskCacheSpeedMode
	c := ioLifeCycleTest(&config, t, w)
	os.Remove(c.DiskPath())
	config.ValueBehaviour = nil
}
func TestAgainstMemAndDiskCacheSpeedModeMultipleCallOnLoadAndOnExpireFuncs(t *testing.T) {
	config := UCacheConfig{
		DiskStorePath: "testCache",
		ExtendExpiry:  false,
		DefaultExpiry: 30 * time.Second,
	}
	w := &wrapper{}
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnLoadFunc: w.OnLoad, OnExpireFunc: w.OnExpire}
	config.StoreType = MemAndDiskCacheSpeedMode
	c := ioLifeCycleMultithreadedGetCallTest(&config, t, w)
	os.Remove(c.DiskPath())
}

func _TestAgainstLongRunningGCDiskCache(t *testing.T) {
	config := UCacheConfig{
		DiskStorePath: "testCache",
		ExtendExpiry:  false,
		DefaultExpiry: 30 * time.Second,
		Verbose:       true,
	}
	w := &wrapper{}
	config.ValueBehaviour = &ValueBehaviour{ValueBeType: CallOnLoadAndOnExpire, OnLoadFunc: w.OnLoad, OnExpireFunc: w.OnExpire}
	config.StoreType = MemAndDiskCacheSpeedMode
	c := ioLifeCycleMultithreadedGetCallTest2(&config, t, w)
	os.Remove(c.DiskPath())
}

type nameIf interface {
	FirstName() string
}
type MyStr struct {
	Name           string
	invisible      int
	Age            int
	Timestamp      time.Time
	onLoadCalled   bool
	onExpireCalled bool
	onClose        bool
}

func (s *MyStr) FirstName() string {
	return s.Name
}
func (s *MyStr) OnLoad() {
	s.onLoadCalled = true
}
func (s *MyStr) OnExpire() {
	s.onExpireCalled = true
}
func (s *MyStr) Close() {
	s.onClose = true
}

type MyStr2 struct {
	Name           string
	invisible      int
	Age            int
	Timestamp      time.Time
	onLoadCalled   bool
	onExpireCalled bool
	onClose        bool
}

func (s *MyStr2) FirstName() string {
	return s.Name
}
func (s *MyStr2) OnLoad() {
	s.onLoadCalled = true
}
func (s *MyStr2) OnExpire() {
	s.onExpireCalled = true
}
func (s *MyStr2) Close() {
	s.onClose = true
}

func ioTest(conf *UCacheConfig, t *testing.T) *UCache {
	c, err := NewUCache(*conf)
	if err != nil {
		t.Error(err)
	}
	c.Put("myStruct", &MyStr{Name: "Artan", Age: 30, Timestamp: time.Now(), invisible: 23})
	c.Put("myStruct2", &MyStr2{Name: "Artan", Age: 30, Timestamp: time.Now(), invisible: 23})
	c.Put("myKey", "my Value")
	c.PutWithOtherExpiry("myKey2", "my Value2", 180*time.Millisecond)
	c.PutWithOtherExpiry("123", 456, 140*time.Millisecond)
	c.Put("1.456", 7.89)

	var myStruct *MyStr
	err = c.Get("myStruct", &myStruct)
	if err != nil || myStruct.Name != "Artan" {
		t.Error(err, myStruct)
	}
	myStruct2 := &MyStr2{}
	err = c.Get("myStruct2", &myStruct2)
	if err != nil || myStruct2.Name != "Artan" {
		t.Error(err, myStruct2)
	}
	var myVal string
	err = c.Get("myKey", &myVal)
	if err != nil {
		t.Error(err, myVal)
	}
	if myVal == "" {
		t.Error("myVal shouldn't be empty", myVal)
	}
	myVal = ""
	time.Sleep(80 * time.Millisecond)
	err = c.Get("myKey", &myVal)
	if err == nil {
		t.Error(err, myVal)
	}
	if myVal != "" {
		t.Error("myVal should be empty because of timeout", myVal)
	}
	myVal = ""
	err = c.Get("myKey2", &myVal)
	if err != nil {
		t.Error(err, myVal)
	}
	if myVal != "my Value2" {
		t.Error("myVal2 shouldn't be empty")
	}
	var intval int
	err = c.Get("123", &intval)
	if err != nil {
		t.Error(err, intval)
	}
	if intval != 456 {
		t.Error("intval shouldn't be empty", intval)
	}

	time.Sleep(time.Millisecond * 280)
	if i := c.Size(); i > 0 {
		t.Error("cache not cleaned properly", i)
	}
	c.Close()
	return c
}

type wrapper struct {
	onLoad      bool
	onExpire    bool
	onLoadCount int
	Names       []string
}

func (me *wrapper) OnLoad(key string, val interface{}) {
	me.onLoad = true
	me.onLoadCount++
	if me.Names == nil {
		me.Names = []string{}
	}
	if n, ok := val.(nameIf); ok {
		me.Names = append(me.Names, n.FirstName())
	}
}

func (me *wrapper) OnExpire(key string, val interface{}) {
	me.onExpire = true
}

func ioLifeCycleTest(conf *UCacheConfig, t *testing.T, w *wrapper) *UCache {
	c, err := NewUCache(*conf)
	if err != nil {
		t.Error(err)
	}
	a := &MyStr{Name: "Artan", Age: 30, Timestamp: time.Now(), invisible: 23}
	c.Put("myStruct", a)
	b := &MyStr2{Name: "Artan", Age: 30, Timestamp: time.Now(), invisible: 23}
	c.Put("myStruct2", b)

	if conf.ValueBehaviour.callOnExpire() && conf.StoreType != DiskCache && conf.StoreType != DiskCacheSpeedMode {
		time.Sleep(time.Millisecond * 200)
		if w != nil && conf.ValueBehaviour.OnExpireFunc != nil {
			if !w.onExpire {
				t.Error("OnExpire not called", conf.StoreType)
			}
		}
		if conf.ValueBehaviour.OnExpireMethodName == "Close" {
			if !a.onClose {
				t.Error("Close not called", conf.StoreType)
			}
			if !b.onClose {
				t.Error("Close not called", conf.StoreType)
			}
		} else if conf.ValueBehaviour.OnExpireFunc == nil {
			if !a.onExpireCalled {
				t.Error("OnExpire not called", conf.StoreType)
			}
			if !b.onExpireCalled {
				t.Error("OnExpire not called", conf.StoreType)
			}
		}
	}

	if conf.StoreType != MemCache {
		a = &MyStr{Name: "Artan", Age: 30, Timestamp: time.Now(), invisible: 23}
		err = c.PutWithOtherExpiry("myStruct", a, 180*time.Millisecond)
		if err != nil {
			t.Error(err)
		}
		b = &MyStr2{Name: "Artan", Age: 30, Timestamp: time.Now(), invisible: 23}
		err = c.PutWithOtherExpiry("myStruct2", b, 180*time.Millisecond)
		if err != nil {
			t.Error(err)
		}
		c.Close()
		c, err = NewUCache(*conf)
		if err != nil {
			t.Error(err)
		}
		if w != nil {
			w.onLoad = false
		}

		a1 := &MyStr{}
		err := c.Get("myStruct", &a1)
		if err != nil || a1.Name != "Artan" {
			t.Error(err, a1)
		}
		if a1.Name != "Artan" {
			t.Error("not the expected object", a1, conf.StoreType)
		}
		if conf.ValueBehaviour.callOnLoad() {
			if w != nil && conf.ValueBehaviour.OnLoadFunc != nil {
				if !w.onLoad || len(w.Names) != 1 || w.Names[0] != "Artan" {
					t.Error("OnLoad not called", conf.StoreType)
				}
			} else if !a1.onLoadCalled {
				t.Error("OnLoad not called", conf.StoreType)
			}
		}
		if w != nil {
			w.onLoad = false
		}
		b = &MyStr2{}
		err = c.Get("myStruct2", &b)
		if err != nil || b.Name != "Artan" {
			t.Error(err)
		}
		if b.Name != "Artan" {
			t.Error("not the expected object", b)
		}
		if conf.ValueBehaviour.callOnLoad() {
			if w != nil && conf.ValueBehaviour.OnLoadFunc != nil {
				if !w.onLoad || len(w.Names) != 2 || w.Names[0] != "Artan" {
					t.Error("OnLoad not called", conf.StoreType, w.Names)
				}
			} else if !b.onLoadCalled {
				t.Error("OnLoad not called", conf.StoreType)
			}
		}
		b.onLoadCalled = false
		if w != nil {
			w.onLoad = false
		}
		var cc *MyStr2
		err = c.Get("myStruct2", &cc)
		if err != nil || cc.Name != "Artan" {
			t.Error(err)
		}
		if cc.Name != "Artan" {
			t.Error("not the expected object")
		}
		if conf.ValueBehaviour.callOnLoad() {
			if w != nil && conf.ValueBehaviour.OnLoadFunc != nil {
				if conf.StoreType == DiskCacheSpeedMode || conf.StoreType == DiskCache {
					if !w.onLoad {
						t.Error("OnLoad not called", conf.StoreType)
					}
				} else {
					if w.onLoad {
						t.Error("OnLoad already called by -> err = c.Get('myStruct2', &b)", conf.StoreType)
					}
				}

			} else {
				if conf.StoreType == DiskCacheSpeedMode || conf.StoreType == DiskCache {
					if !cc.onLoadCalled {
						t.Error("OnLoad not called", conf.StoreType)
					}
				} else {
					if cc.onLoadCalled {
						t.Error("OnLoad already called by -> err = c.Get('myStruct2', &b)", conf.StoreType)
					}
				}

			}
		}
	}
	time.Sleep(240 * time.Millisecond)
	if i := c.Size(); i > 0 {
		t.Error("cache not cleaned properly", i, c.storeType)
	}
	c.Close()
	return c
}

func ioLifeCycleMultithreadedGetCallTest(conf *UCacheConfig, t *testing.T, w *wrapper) *UCache {
	c, err := NewUCache(*conf)
	if err != nil {
		t.Error(err)
	}
	a := &MyStr{Name: "Artan", Age: 30, Timestamp: time.Now(), invisible: 23}
	c.Put("myStruct", a)
	c.Close()
	//test more than 10 times as it is a multi threaded test but it still doesn't guaranty it is correct
	for count := 0; count < 15; count++ {
		w.onLoadCount = 0
		w.onLoad = false
		c, err = NewUCache(*conf)
		if err != nil {
			t.Error(err)
		}
		wait := true
		for i := 0; i < 6; i++ {
			go func() {
				for wait {
					time.Sleep(time.Nanosecond * 1)
				}
				var a1 *MyStr
				err = c.Get("myStruct", &a1)
				if err != nil {
					t.Error(err)
				}
				if a1 == nil {
					t.Error("a1 is nil")
				}
			}()
		}
		time.Sleep(time.Millisecond * 10)
		wait = false

		time.Sleep(time.Millisecond * 500)
		if !w.onLoad || w.onLoadCount != 1 {
			t.Error("OnLoad not called", w.onLoadCount)
		}
		c.Close()
	}
	return c
}

func ioLifeCycleMultithreadedGetCallTest2(conf *UCacheConfig, t *testing.T, w *wrapper) *UCache {
	c, err := NewUCache(*conf)
	if err != nil {
		t.Error(err)
	}
	a := &MyStr{Name: "Artan", Age: 30, Timestamp: time.Now(), invisible: 23}
	c.Put("myStruct", a)
	c.Close()
	//test more than 10 times as it is a multi threaded test but it still doesn't guaranty it is correct
	for count := 0; count < 15; count++ {
		w.onLoadCount = 0
		w.onLoad = false
		c, err = NewUCache(*conf)
		if err != nil {
			t.Error(err)
		}
		wait := true
		for i := 0; i < 6; i++ {
			go func() {
				for wait {
					time.Sleep(time.Nanosecond * 1)
				}
				id := i
				for ii := 0; ii < 1000000; ii++ {
					err = c.Put(fmt.Sprintf("myStruct %d %d", ii, id), "asfasdf")
					if err != nil {
						t.Error(err)
					}
				}

			}()
		}
		time.Sleep(time.Millisecond * 10)
		wait = false

		time.Sleep(time.Minute * 1)
		if !w.onLoad || w.onLoadCount != 1 {
			t.Error("OnLoad not called", w.onLoadCount)
		}
		c.Close()
	}
	return c
}

var benchConfig = UCacheConfig{
	DiskStorePath: "testCache",
	DefaultExpiry: 30 * time.Second,
}

func BenchmarkMemCache(b *testing.B) {
	benchPut(MemCache, b)
}
func BenchmarkDiskCache(b *testing.B) {
	benchPut(DiskCache, b)
}
func BenchmarkMemAndDiskCache(b *testing.B) {
	benchPut(MemAndDiskCache, b)
}
func BenchmarkMemAndDiskCacheSpeedMode(b *testing.B) {
	benchPut(MemAndDiskCacheSpeedMode, b)
}
func BenchmarkDiskCacheSpeedMode(b *testing.B) {
	benchPut(DiskCacheSpeedMode, b)
}

func benchPut(st StoreType, b *testing.B) {
	benchConfig.StoreType = st
	c, err := NewUCache(benchConfig)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		c.PutWithOtherExpiry(fmt.Sprintf("myKey2%v", i), MyStr{Name: "Artan", Age: 30}, 80*time.Second)
	}
	c.Close()
	if st != MemCache {
		os.Remove(c.DiskPath())
	}
}
