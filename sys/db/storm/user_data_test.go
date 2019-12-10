package storm

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ProxeusApp/proxeus-core/sys/file"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

func TestPutGetData(t *testing.T) {
	db, err := NewUserDataDB("./")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		db.Close()
		db.Remove()
	}()
	u := &model.User{ID: "123", Role: model.ROOT}
	usrData := &model.UserDataItem{Name: "some name", Detail: "detail..", Data: map[string]interface{}{"input": map[string]interface{}{"abcField": 123}, "123": "along input"}}
	err = db.Put(u, usrData)
	if err != nil {
		t.Error(err)
	}
	data := map[string]interface{}{"input": map[string]interface{}{"list": []interface{}{1, 2, 3, 4, "hello"}, "someInt": 1234}}
	err = db.PutData(u, usrData.ID, data)
	if err != nil {
		t.Error(err)
	}
	newUsrData, err := db.Get(u, usrData.ID)
	if err != nil {
		t.Error(err)
	}
	if newUsrData.Data == nil {
		t.Error("data is nil")
	}
	if innerMap, ok := newUsrData.Data["input"].(map[string]interface{}); ok {
		if someInt, ok := innerMap["someInt"].(int64); ok {
			if someInt != 1234 {
				t.Error("someInt missing", someInt)
			}
		} else {
			t.Error("someInt missing")
		}
		if list, ok := innerMap["list"].([]interface{}); ok {
			if i, ok := list[0].(int64); ok && i != 1 {
				t.Error("not 1")
			}
			if i, ok := list[1].(int64); ok && i != 2 {
				t.Error("not 2")
			}
			if i, ok := list[2].(int64); ok && i != 3 {
				t.Error("not 3")
			}
			if i, ok := list[3].(int64); ok && i != 4 {
				t.Error("not 4")
			}
			if i, ok := list[4].(string); ok && i != "hello" {
				t.Error("not hello")
			}
		}
	} else {
		t.Error("input missing")
	}

	data = map[string]interface{}{"input": []interface{}{1, 2, 3, 4, "hello"}}

	err = db.PutData(u, usrData.ID, data)
	if err != nil {
		t.Error(err)
	}
	newUsrData, err = db.Get(u, usrData.ID)
	if err != nil {
		t.Error(err)
	}
	if newUsrData.Data == nil {
		t.Error("data is nil")
	}

	if list, ok := newUsrData.Data["input"].([]interface{}); ok {
		if i, ok := list[0].(uint16); ok && i != 1 {
			t.Error("not 1")
		}
		if i, ok := list[1].(uint16); ok && i != 2 {
			t.Error("not 2")
		}
		if i, ok := list[2].(uint16); ok && i != 3 {
			t.Error("not 3")
		}
		if i, ok := list[3].(uint16); ok && i != 4 {
			t.Error("not 4")
		}
		if i, ok := list[4].(string); ok && i != "hello" {
			t.Error("not hello")
		}
	}

	//file input

	fi := db.NewFile(u, file.Meta{Name: "abc", ContentType: "ct", Size: 123})
	db.PutData(u, usrData.ID, map[string]interface{}{"input": fi})

	newUsrData, err = db.Get(u, usrData.ID)
}

func TestPutGetDataFile(t *testing.T) {
	dir := "./testDir"
	db, err := NewUserDataDB(dir)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		db.Close()
		db.Remove()
	}()
	u := &model.User{ID: "123", Role: model.ROOT}
	usrData := &model.UserDataItem{Name: "some name", Detail: "detail..", Data: map[string]interface{}{"input": map[string]interface{}{"abcField": 123}, "123": "along input"}}
	err = db.Put(u, usrData)
	if err != nil {
		t.Error(err)
	}

	fileData := "hello"
	err = ioutil.WriteFile(filepath.Join(db.baseFilePath, "abc"), []byte(fileData), 0600)
	if err != nil {
		t.Error(err)
		return
	}
	fi := db.NewFile(u, file.Meta{Name: "MyFile"})
	f, err := os.Open(filepath.Join(db.baseFilePath, "abc"))
	if err != nil {
		t.Error(err)
		return
	}
	n, err := fi.Write(f)
	if err != nil {
		t.Error(n, err)
		return
	}
	err = db.PutData(u, usrData.ID, map[string]interface{}{"file": fi})
	if err != nil {
		t.Error(n, err)
		return
	}
	shouldBeSameFi, err := db.GetDataFile(u, usrData.ID, "file")
	if err != nil {
		t.Error(n, err)
		return
	}
	buf := &bytes.Buffer{}
	n, err = shouldBeSameFi.Read(buf)
	if err != nil {
		t.Error(n, err)
		return
	}
	if buf.String() != fileData {
		t.Error(n, err)
	}

	//put another file
	fi = db.NewFile(u, file.Meta{Name: "MyFile2"})
	n, err = fi.Write(f)
	if err != nil {
		t.Error(n, err)
		return
	}
	err = db.PutData(u, usrData.ID, map[string]interface{}{"file2": fi})
	if err != nil {
		t.Error(err)
		return
	}
	fi, err = db.GetDataFile(u, usrData.ID, "file2")
	if err != nil {
		t.Error(err)
		return
	}
	n, err = fi.Read(buf)
	if err != nil {
		t.Error(n, err)
		return
	}
	if buf.String() != fileData {
		t.Error(n, err)
	}

	usrItem, err := db.Get(u, usrData.ID)
	if err != nil {
		log.Println(err)
	}
	log.Println("DATA", usrItem.Data)

	db.Close()
	//os.RemoveAll(dir)
	//log.Println("err", usrData)
}

func TestGetAssignedUsers(t *testing.T) {
	//m := make(map[string]*model.User)
	//m["111"] = &model.User{ID: "111"}

	hello("123", "asfd")
}

func hello(id ...string) {
	specificIds := len(id) > 0
	log.Println(specificIds)
	s := makeSimpleQuery(map[string]interface{}{"include": id})
	//testMethod(m)
	log.Println(s.include)
}

func testMethod(m map[string]*model.User) {
	m["123"] = &model.User{ID: "123"}
}
