package file

import (
	"fmt"
	"log"
	"regexp"
	"testing"
)

func TestMapIO_MergeWith(t *testing.T) {
	mapIO := MapIO{}
	mapIO.MergeWith(map[string]interface{}{
		"input": map[string]interface{}{
			"abcField": 123,
		},
		"123": "along input",
	})
	data := map[string]interface{}{
		"input": map[string]interface{}{
			"list":    []interface{}{1, 2, 3, 4, "hello"},
			"someInt": 1234,
		},
	}
	mapIO.MergeWith(data)
	if innerMap, ok := mapIO["input"].(map[string]interface{}); ok {
		if someInt, ok := innerMap["someInt"].(int); ok {
			if someInt != 1234 {
				t.Error("someInt missing", someInt)
			}
		} else {
			t.Error("someInt missing")
		}
		if list, ok := innerMap["list"].([]interface{}); ok {
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
	} else {
		t.Error("input missing")
	}

	data = map[string]interface{}{"input": []interface{}{1, 2, 3, 4, "hello"}}
	mapIO.MergeWith(data)

	if list, ok := mapIO["input"].([]interface{}); ok {
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
}

type ABC struct {
	MyField  map[string]interface{}
	myField2 *ABC
	Name     string
}

func TestMapIO_Get(t *testing.T) {
	mapIO := MapIO{}
	mapIO.MergeWith(map[string]interface{}{
		"input": map[string]interface{}{
			"abcField": 123,
		},
		"123": "along input",
	})
	data := map[string]interface{}{
		"input": map[string]interface{}{
			"list":    []interface{}{1, 2, 3, 4, "hello", New("", Meta{Name: "a", ContentType: "b", Size: 1})},
			"someInt": 1234,
		},
		"list": []interface{}{
			1,
			2,
			map[string]interface{}{"size": 1, "path": "123", "contentType": "ct", "name": "ab"}, 4, "hello", New("", Meta{Name: "a", ContentType: "b", Size: 1})},
		"ABC": ABC{Name: "MyName", MyField: map[string]interface{}{"Inner": 111}},
	}

	mapIO.MergeWith(data)
	intrfc := mapIO.Get("ABC.MyField.Inner")
	if nr, ok := intrfc.(int); !ok || nr != 111 {
		t.Error(nr, intrfc)
	}
	intrfc = mapIO.Get("input.someInt")
	if nr, ok := intrfc.(int); !ok || nr != 1234 {
		t.Error(nr, intrfc)
	}
	intrfc = mapIO.Get("ABC.Name")
	if nr, ok := intrfc.(string); !ok || nr != "MyName" {
		t.Error(nr, intrfc)
	}
	intrfc = mapIO.Get("input.list[1]")
	if nr, ok := intrfc.(int); !ok || nr != 2 {
		t.Error(nr)
	}
	intrfc = mapIO.Get("input.list[4]")
	if nr, ok := intrfc.(string); !ok || nr != "hello" {
		t.Error(nr)
	}
	fio := mapIO.GetFileInfo("", "input.list[5]")
	if fio == nil {
		t.Error(fio)
	}
	fio = mapIO.GetFileInfo("", "list[5]")
	if fio == nil {
		t.Error(fio)
	}
	fio = mapIO.GetFileInfo("", "list[2]")
	if fio == nil {
		t.Error(fio)
	}
	fio = mapIO.GetFileInfo("", "list[12]")
	if fio != nil {
		t.Error(fio)
	}
	fio = mapIO.GetFileInfo("", "list[6]")
	if fio != nil {
		t.Error(fio)
	}
	fio = mapIO.GetFileInfo("", "list[-1]")
	if fio != nil {
		t.Error(fio)
	}
}

func TestMapIO_MakeFileInfos(t *testing.T) {
	mapIO := MapIO{}
	mapIO["file"] = map[string]interface{}{"path": "123", "size": 123, "contentType": "ct", "name": "abc"}
	mapIO.MergeWith(map[string]interface{}{
		"input": map[string]interface{}{
			"abcField":    123,
			"anotherFile": map[string]interface{}{"path": "123", "size": 123, "contentType": "ct", "name": "abc"},
		},
		"123": "along input",
	})
	mapIO.MakeFileInfos("abc")
	shouldBeFi := mapIO.Get("file")
	if fi, ok := shouldBeFi.(*IO); !ok {
		t.Error(fi)
	}
	shouldBeFi = mapIO.Get("input.anotherFile")
	if fi, ok := shouldBeFi.(*IO); !ok {
		t.Error(fi)
	}
}

func TestMapIO_EmptyKey(t *testing.T) {
	mapIO := MapIO{}
	mapIO.MergeWith(map[string]interface{}{
		"":     "123",
		"name": "abc",
	})
	if len(fmt.Sprint(mapIO)) != len(`map[name:abc]`) {
		t.Error("wrong size")
	}
}

func TestMapIO_Get2(t *testing.T) {
	var re = regexp.MustCompile(`\[|\]|\.`)
	var str = `abc.omfg[1].abcd[4].ad`
	arr := re.Split(str, -1)
	log.Println(arr)
}
