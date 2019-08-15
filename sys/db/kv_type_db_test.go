package db

import (
	"testing"
)

type complexStruct struct {
	Bla  int
	Blup float64
	Omfg bool
	Rofl string
}

func Test_KVTypeMEMGet(t *testing.T) {
	ss, err := NewKVStore(nil, "")
	var first *complexStruct
	first = &complexStruct{Bla: 1, Omfg: true, Rofl: "yes"}
	err = ss.Put("first", first)
	if err != nil {
		t.Error(err)
		return
	}
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
	err = ss.Put("myInt", myInt)
	if err != nil {
		t.Error(err)
		return
	}
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
	err = ss.Put("myFloat", myFloat)
	if err != nil {
		t.Error(err)
		return
	}
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
}
