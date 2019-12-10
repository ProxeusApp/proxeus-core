package session

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type complexStruct struct {
	Bla  int
	Blup float64
	Omfg bool
	Rofl string
}

func TestJSONSanityCheck(t *testing.T) {
	s := &Session{id: "myID", userID: "myUserID", rights: model.ADMIN, userName: "myUserName", sessionDir: "mySessionDir"}
	bts, err := json.Marshal(s)
	if err != nil {
		t.Error(err, string(bts))
	}
	s2 := &Session{}
	err = json.Unmarshal(bts, &s2)
	if err != nil {
		t.Error(err)
	}
}

func TestSessionMemoryStorage(t *testing.T) {
	var err error
	mainDir := "testSessionDir"
	m, er := NewManager(mainDir, time.Second*1)
	if er != nil {
		t.Error(er)
	}
	var ss *Session
	ss, err = m.New("1", "name", model.CREATOR)
	var first *complexStruct
	first = &complexStruct{Bla: 1, Omfg: true, Rofl: "yes"}
	ss.Put("first", first)
	first = nil
	err = ss.Get("first", &first)
	if err != nil {
		t.Error(err)
		return
	}
	if first == nil || first.Bla != 1 || !first.Omfg || first.Rofl != "yes" {
		t.Error(err)
		return
	}
	var myInt uint = 1232323
	ss.Put("myInt", myInt)
	myInt = 28383
	err = ss.Get("myInt", &myInt)
	if err != nil {
		t.Error(err)
		return
	}
	if myInt != 1232323 {
		t.Error(err)
		return
	}
	var myFloat = 1232.323
	ss.Put("myFloat", myFloat)
	myFloat = 283.383
	err = ss.Get("myFloat", &myFloat)
	if err != nil {
		t.Error(err)
		return
	}
	if myFloat != 1232.323 {
		t.Error(err)
		return
	}
	ss.Kill()
	m.Close()
	os.RemoveAll(mainDir)
}
