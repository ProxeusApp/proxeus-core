package sys

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var dbconfig = map[string]interface{}{"absolute.key.path": false, "sort": map[string]interface{}{"firstObject": "asc"}}

func TestWriteRead(t *testing.T) {
	x := struct {
		Foo string
		Bar int
	}{"foo", 2}

	v := reflect.ValueOf(x)

	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	fmt.Println(values)
}

func TestRemWildcards(t *testing.T) {
	var dbPath = "TestWriteReadTypeConversion.db"
	db, err := Open(dbPath)
	fmt.Println(string(db.wildcardPieces("alskdjf*sdjf.lsjfskdfj*")))
	fmt.Println(string(db.wildcardPieces("alskdjf***sdjf.lsjfskdfj*")))
	err = os.Remove(dbPath)
	if err != nil {
		t.Error("cannot delete db file ", err)
	}
}

func TestUnifySelector(t *testing.T) {
	var dbPath = "TestWriteReadTypeConversion2.db"
	db, err := Open(dbPath)
	sel, keys := db.unify(`omg.blabla["myemail@gmail.com"].omg.blabla[1]`)
	fmt.Println(string(*sel))
	for _, b := range *keys {
		fmt.Println(string(b))
	}
	err = os.Remove(dbPath)
	if err != nil {
		t.Error("cannot delete db file ", err)
	}
}

func TestAbsoluteKeyPath(t *testing.T) {
	var dbPath = "TestWriteReadTypeConversion3.db"
	db, err := Open(dbPath)
	var config = map[string]interface{}{"absolute.key.path": true, "sort": map[string]interface{}{"firstObject": "asc"}}
	_ = db.Write(map[string]interface{}{
		"firstObject": map[string]interface{}{
			"av@futuretek.ch": "abc",
			"boolean":         true,
			"string int":      "123456",
			"mayArray": map[string]interface{}{
				"mayArrayInner1": []interface{}{
					map[string]interface{}{
						"mayArrayInner3.kjkjsdf": map[string]interface{}{
							"mayArrayInner4": map[string]interface{}{
								"a": "abc",
								"b": "asd",
							},
						},
					},
				},
			},
		}}, config)
	res := db.Read([]string{"firstObject.mayArray.mayArrayInner1[0].mayArrayInner3.kjkjsdf.mayArrayInner4"}, config)
	myMap, ok := res[0].(map[string]interface{})
	if !ok {
		t.Error("entry is not a map[string]interface{} as expected!")
	}
	fancyPrint(myMap)
	err = os.Remove(dbPath)
	if err != nil {
		t.Error("cannot delete db file ", err)
	}
}

func TestArray(t *testing.T) {
	a := make([]int, 0)
	for i := 0; i < 100; i++ {
		a = append(a, i+1)
	}
	fmt.Println(a)
}

func TestDelete(t *testing.T) {
	var dbPath = "TestWriteReadTypeConversion4.db"
	db, err := Open(dbPath)
	var i int = 2147483647
	_ = db.Write(map[string]interface{}{
		"firstObject": map[string]interface{}{
			"av@futuretek.ch": i,
			"boolean":         true,
			"string int":      "123456",
			"mayArray": []interface{}{
				"a",
				"b",
			},
		}}, dbconfig)
	fmt.Println("dump 1-------------------------")
	db.Dump(nil)
	fmt.Println("dump 1-------------------------")
	db.Delete([]string{"firstObject.string int"}, nil)
	fmt.Println("dump 2-------------------------")
	db.Dump(nil)
	fmt.Println("dump 2-------------------------")

	err = os.Remove(dbPath)
	if err != nil {
		t.Error("cannot delete db file ", err)
	}
}

func TestWriteReadTypeConversion(t *testing.T) {
	var dbPath = "TestWriteReadTypeConversion5.db"
	db, err := Open(dbPath)
	var i int = 2147483647
	var imin int = -2147483648
	var i8 int8 = 127
	var i8min int8 = -128
	var i82 byte = 255
	var i82min byte = 0
	var i16 int16 = 32767
	var i16min int16 = -32768
	var i32 int32 = 2147483647
	var i32min int32 = -2147483648
	var i64 int64 = 9223372036854775807
	var i64min int64 = 9223372036854775807
	var ui uint = 2147483647
	var uimin uint = 0
	var ui8 uint8 = 100
	var ui8min uint8 = 0
	var ui16 uint16 = 65535
	var ui16min uint16 = 0
	var ui32 uint32 = 4294967295
	var ui32min uint32 = 0
	var ui64 uint64 = 18446744073709551615
	var ui64min uint64 = 0
	var f32 float32 = 123456.123456
	var f32min float32 = -123456.123456
	var f64 = 123456789101112.12345678910
	var f64min = -12345678910.12345678910
	//var c64 complex64 = 2.323423
	//var c128 complex128 = 12345678910.12345678910
	var b = true
	_ = db.Write(map[string]interface{}{
		"firstObject": map[string]interface{}{
			"av@futuretek.ch": i,
			"boolean":         true,
			"string int":      "123456",
			"i":               i,
			"imin":            imin,
			"i8":              i8,
			"i8min":           i8min,
			"i82":             i82,
			"i82min":          i82min,
			"i16":             i16,
			"i16min":          i16min,
			"i32":             i32,
			"i32min":          i32min,
			"i64":             i64,
			"i64min":          i64min,
			"f32":             f32,
			"f32min":          f32min,
			"f64":             f64,
			"f64min":          f64min,
			"ui":              ui,
			"uimin":           uimin,
			"ui8":             ui8,
			"ui8min":          ui8min,
			"ui16":            ui16,
			"ui16min":         ui16min,
			"ui32":            ui32,
			"ui32min":         ui32min,
			"ui64":            ui64,
			"ui64min":         ui64min,
			//"c64":c64,
			//"c128":c128,
			"account": map[string]interface{}{
				"i":    i,
				"i8":   i8,
				"i16":  i16,
				"i32":  i32,
				"i64":  i64,
				"f32":  f32,
				"f64":  f64,
				"ui":   ui,
				"ui16": ui16,
				"ui32": ui32,
				"ui64": ui64,
			},
			"b": b,
			"mayArray": []interface{}{
				"a",
				"b",
				i8,
			},
			"mayArray2": []string{
				"a",
				"b",
			},
			"mayArray3": []int{
				4,
				2,
			},
		}}, dbconfig)
	res := db.Read([]string{"firstObject"}, dbconfig)
	if len(res) == 0 {
		t.Error("read length is 0! should be 1")
	}

	myMap, ok := res[0].(map[string]interface{})
	if !ok {
		t.Error("entry is not a map[string]interface{} as expected!")
	}
	if ok, _ = dbconfig["absolute.key.path"].(bool); ok {
		myMap = myMap["firstObject"].(map[string]interface{})
	}

	//fancyPrint(myMap)
	db.Dump(nil)
	err = os.Remove(dbPath)
	if err != nil {
		t.Error("cannot delete db file ", err)
	}
	var val interface{}
	var tmpObj interface{}
	val = myMap["i"]
	tmpObj, ok = val.(int)
	if val != i || !ok {
		t.Error("value not as expected", val, i, ok)
	}
	val = myMap["imin"]
	tmpObj, ok = val.(int)
	if val != imin || !ok {
		t.Error("value not as expected", val, imin, ok)
	}
	val = myMap["i8"]
	tmpObj, ok = val.(int8)
	if val != i8 || !ok {
		t.Error("value not as expected", val, i8)
	}
	val = myMap["i8min"]
	tmpObj, ok = val.(int8)
	if val != i8min || !ok {
		t.Error("value not as expected", val, i8min)
	}
	val = myMap["i82"]
	tmpObj, ok = val.(byte)
	if val != i82 || !ok {
		t.Error("value not as expected", val, i8)
	}
	val = myMap["i82min"]
	tmpObj, ok = val.(byte)
	if val != i82min || !ok {
		t.Error("value not as expected", val, i8min)
	}
	val = myMap["i16"]
	tmpObj, ok = val.(int16)
	if val != i16 || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["i16min"]
	tmpObj, ok = val.(int16)
	if val != i16min || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["i32"]
	tmpObj, ok = val.(int32)
	if val != i32 || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["i32min"]
	tmpObj, ok = val.(int32)
	if val != i32min || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["i64"]
	tmpObj, ok = val.(int64)
	if val != i64 || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["i64min"]
	tmpObj, ok = val.(int64)
	if val != i64min || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["ui"]
	tmpObj, ok = val.(uint)
	if val != ui || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["uimin"]
	tmpObj, ok = val.(uint)
	if val != uimin || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["ui8"]
	tmpObj, ok = val.(uint8)
	if val != ui8 || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["ui8min"]
	tmpObj, ok = val.(uint8)
	if val != ui8min || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["ui16"]
	tmpObj, ok = val.(uint16)
	if val != ui16 || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["ui16min"]
	tmpObj, ok = val.(uint16)
	if val != ui16min || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["ui32"]
	tmpObj, ok = val.(uint32)
	if val != ui32 || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["ui32min"]
	tmpObj, ok = val.(uint32)
	if val != ui32min || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["ui64"]
	tmpObj, ok = val.(uint64)
	if val != ui64 || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["ui64min"]
	tmpObj, ok = val.(uint64)
	if val != ui64min || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["f32"]
	tmpObj, ok = val.(float32)
	if val != f32 || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["f32min"]
	tmpObj, ok = val.(float32)
	if val != f32min || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["f64"]
	tmpObj, ok = val.(float64)
	if val != f64 || !ok {
		t.Error("value not as expected", val)
	}
	val = myMap["f64min"]
	tmpObj, ok = val.(float64)
	if val != f64min || !ok {
		t.Error("value not as expected", val)
	}
	if tmpObj == nil {
		fmt.Println(tmpObj)
	}
}

func fancyPrint(d interface{}) {
	_fancyPrint(d, 0, "")
}

func _fancyPrint(d interface{}, intent int, key string) {
	switch v := d.(type) {
	case []interface{}:
		addIntent(intent)
		fmt.Print("[\n")
		for i := 0; i < len(v); i++ {
			_fancyPrint(v[i], intent+1, "")
		}
		addIntent(intent)
		fmt.Print("]\n")
	case map[string]interface{}:
		addIntent(intent)
		if key == "" {
			fmt.Print("{\n")
		} else {
			fmt.Printf("%s: {\n", key)
		}
		for mk := range v {
			_fancyPrint(v[mk], intent+1, mk)
		}
		addIntent(intent)
		fmt.Print("}\n")
	default:
		addIntent(intent)
		fmt.Printf("%v %v\n", reflect.TypeOf(v), v)
	}
}

func addIntent(size int) {
	for i := 0; i < size; i++ {
		fmt.Print("\t")
	}
}
