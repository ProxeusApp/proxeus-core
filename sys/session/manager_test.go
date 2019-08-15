package session

import (
	"os"
	"testing"
	"time"

	"git.proxeus.com/core/central/sys/model"
)

type myNotify struct {
	myOnCreatedMap map[string]bool
	myOnLoadMap    map[string]bool
	myOnExpireMap  map[string]bool
	myOnRemovedMap map[string]bool
}

func (me *myNotify) OnSessionCreated(key string, s *Session) {
	me.myOnCreatedMap[key] = true
}
func (me *myNotify) OnSessionLoaded(key string, s *Session) {
	me.myOnLoadMap[key] = true
}
func (me *myNotify) OnSessionExpired(key string, s *Session) {
	me.myOnExpireMap[key] = true
}
func (me *myNotify) OnSessionRemoved(key string) {
	me.myOnRemovedMap[key] = true
}

func TestOnCreatedOnLoadOnExpireOnRemove(t *testing.T) {
	sessDir := "./testSessionDir"
	expiry := time.Millisecond * 800
	myOnCreatedMap := make(map[string]bool)
	myOnLoadMap := make(map[string]bool)
	myOnExpireMap := make(map[string]bool)
	myOnRemovedMap := make(map[string]bool)
	sm, err := NewManagerWithNotify(sessDir, expiry, &myNotify{
		myOnLoadMap:    myOnLoadMap,
		myOnCreatedMap: myOnCreatedMap,
		myOnExpireMap:  myOnExpireMap,
		myOnRemovedMap: myOnRemovedMap,
	})
	if err != nil {
		t.Error(err)
	}
	s, err := sm.New("1", "abc", model.ADMIN)
	if err != nil {
		t.Error(err)
	}
	if created, exists := myOnCreatedMap[s.ID()]; !exists || !created {
		t.Error(s.ID(), "not created")
	}
	s1ID := s.ID()
	s, err = sm.New("2", "aaaa", model.USER)
	if err != nil {
		t.Error(err)
	}
	if created, exists := myOnCreatedMap[s.ID()]; !exists || !created {
		t.Error(s.ID(), "not created")
	}
	s2ID := s.ID()
	sm.Close()

	sm, err = NewManagerWithNotify(sessDir, expiry, &myNotify{
		myOnLoadMap:    myOnLoadMap,
		myOnCreatedMap: myOnCreatedMap,
		myOnExpireMap:  myOnExpireMap,
		myOnRemovedMap: myOnRemovedMap,
	})
	if err != nil {
		t.Error(err)
	}
	//onLoad must be called for this session
	s, err = sm.Get(s1ID)
	if err != nil {
		t.Error(err)
	}
	if loaded, exists := myOnLoadMap[s.ID()]; !exists || !loaded {
		t.Error(s.ID(), "not loaded")
	}
	if s.manager == nil {
		t.Error("manager shouldn't be nil")
	}
	//onLoad must be called for this session
	s, err = sm.Get(s2ID)
	if err != nil {
		t.Error(err)
	}
	if loaded, exists := myOnLoadMap[s.ID()]; !exists || !loaded {
		t.Error(s.ID(), "not loaded")
	}
	if s.manager == nil {
		t.Error("manager shouldn't be nil")
	}
	time.Sleep(time.Second * 1)
	if exired, exists := myOnExpireMap[s1ID]; !exists || !exired {
		t.Error(s1ID, "not exired")
	}
	if exired, exists := myOnExpireMap[s2ID]; !exists || !exired {
		t.Error(s2ID, "not exired")
	}
	_, err = os.Stat(sm.constructSessionDir(s1ID))
	if !os.IsNotExist(err) {
		t.Error(err)
	}
	_, err = os.Stat(sm.constructSessionDir(s2ID))
	if !os.IsNotExist(err) {
		t.Error(err)
	}
	//insert again and test remove
	s, err = sm.New("1", "abc", model.ADMIN)
	if err != nil {
		t.Error(err)
	}
	s1ID = s.ID()
	sm.Remove(s1ID)
	if removed, exists := myOnRemovedMap[s1ID]; !exists || !removed {
		t.Error(s1ID, "not removed")
	}
	_, err = os.Stat(sm.constructSessionDir(s1ID))
	if !os.IsNotExist(err) {
		t.Error(err)
	}
	s, err = sm.New("2", "aaaa", model.USER)
	if err != nil {
		t.Error(err)
	}
	s2ID = s.ID()
	sm.Remove(s2ID)
	if removed, exists := myOnRemovedMap[s2ID]; !exists || !removed {
		t.Error(s2ID, "not removed")
	}
	_, err = os.Stat(sm.constructSessionDir(s2ID))
	if !os.IsNotExist(err) {
		t.Error(err)
	}
	sm.Close()
	os.RemoveAll(sessDir)
}

func TestExpireWithoutLoading(t *testing.T) {
	sessDir := "./testSessionDir"
	expiry := time.Millisecond * 800
	sm, err := NewManager(sessDir, expiry)
	if err != nil {
		t.Error(err)
	}
	s, err := sm.New("1", "abc", model.ADMIN)
	if err != nil {
		t.Error(err)
	}
	s1ID := s.ID()
	s, err = sm.New("2", "aaaa", model.USER)
	if err != nil {
		t.Error(err)
	}
	s2ID := s.ID()
	sm.Close()
	sm, err = NewManager(sessDir, expiry)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 2)
	_, err = os.Stat(sm.constructSessionDir(s1ID))
	if !os.IsNotExist(err) {
		t.Error(err, sm.constructSessionDir(s1ID))
	}
	_, err = os.Stat(sm.constructSessionDir(s2ID))
	if !os.IsNotExist(err) {
		t.Error(err, sm.constructSessionDir(s2ID))
	}
	sm.Close()
	//os.RemoveAll(sessDir)
}

type MySessionObject struct {
	closeCalled bool
}

func (me *MySessionObject) Close() {
	me.closeCalled = true
}

type MySessionObject2 struct {
	closeCalled bool
}

func (me *MySessionObject2) Close() error {
	me.closeCalled = true
	return nil
}

type MySessionObject3 struct {
	closeCalled bool
}

func (me *MySessionObject3) Close() (string, error) {
	me.closeCalled = true
	return "", nil
}

func TestExtendExpiryAndCloseOnSessionMemStore(t *testing.T) {
	sessDir := "./testSessionDir"
	expiry := time.Millisecond * 800
	myOnCreatedMap := make(map[string]bool)
	myOnLoadMap := make(map[string]bool)
	myOnExpireMap := make(map[string]bool)
	myOnRemovedMap := make(map[string]bool)
	sm, err := NewManagerWithNotify(sessDir, expiry, &myNotify{
		myOnLoadMap:    myOnLoadMap,
		myOnCreatedMap: myOnCreatedMap,
		myOnExpireMap:  myOnExpireMap,
		myOnRemovedMap: myOnRemovedMap,
	})
	if err != nil {
		t.Error(err)
	}
	s, err := sm.New("1", "abc", model.ADMIN)
	if err != nil {
		t.Error(err)
	}
	s1ID := s.ID()

	obj1 := &MySessionObject{}
	s.Put("obj1", obj1)
	obj2 := &MySessionObject2{}
	s.Put("obj2", obj2)
	obj3 := &MySessionObject3{}
	s.Put("obj3", obj3)
	//test if Close call will fail because of other types
	s.Put("int", 1)
	s.Put("string", "hello")
	s.Put("bool", true)
	s.Put("float", 123.123)

	time.Sleep(time.Millisecond * 750)
	s, err = sm.Get(s1ID)
	if err != nil || s == nil {
		t.Error(err, s)
	}
	if exired, exists := myOnExpireMap[s1ID]; exists || exired {
		t.Error(s1ID, "should't exire")
	}
	time.Sleep(time.Millisecond * 750)
	s, err = sm.Get(s1ID)
	if err != nil || s == nil {
		t.Error(err, s)
	}
	if exired, exists := myOnExpireMap[s1ID]; exists || exired {
		t.Error(s1ID, "should't exire")
	}

	time.Sleep(time.Millisecond * 750)
	s, err = sm.Get(s1ID)
	if err != nil || s == nil {
		t.Error(err, s)
	}
	if exired, exists := myOnExpireMap[s1ID]; exists || exired {
		t.Error(s1ID, "should't exire")
	}
	time.Sleep(time.Second * 1)
	if exired, exists := myOnExpireMap[s1ID]; !exists || !exired {
		t.Error(s1ID, "not exired")
	}
	if !obj1.closeCalled {
		t.Error("obj1 close not called")
	}
	if !obj2.closeCalled {
		t.Error("obj2 close not called")
	}
	if !obj3.closeCalled {
		t.Error("obj3 close not called")
	}
	_, err = os.Stat(sm.constructSessionDir(s1ID))
	if !os.IsNotExist(err) {
		t.Error(err)
	}
	sm.Close()
	os.RemoveAll(sessDir)
}

func TestCloseOnSessionMemStoreWhenClosingManager(t *testing.T) {
	sessDir := "./testSessionDir"
	expiry := time.Second * 3
	sm, err := NewManager(sessDir, expiry)
	if err != nil {
		t.Error(err)
	}
	s, err := sm.New("1", "abc", model.ADMIN)
	if err != nil {
		t.Error(err)
	}

	obj1 := &MySessionObject{}
	s.Put("obj1", obj1)
	obj2 := &MySessionObject2{}
	s.Put("obj2", obj2)
	obj3 := &MySessionObject3{}
	s.Put("obj3", obj3)
	//test if Close call will fail because of other types
	s.Put("int", 1)
	s.Put("string", "hello")
	s.Put("bool", true)
	s.Put("float", 123.123)

	sm.Close()

	if !obj1.closeCalled {
		t.Error("obj1 close not called")
	}
	if !obj2.closeCalled {
		t.Error("obj2 close not called")
	}
	if !obj3.closeCalled {
		t.Error("obj3 close not called")
	}

	os.RemoveAll(sessDir)
}
