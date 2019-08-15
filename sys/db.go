package sys

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"

	//"sort"
	"crypto/rand"
)

var (
	defaultOptions = map[string]map[string]interface{}{
		"global": {
			"region": []byte("main"),
		},
		"read": {
			"absolute.key.path":      true, // query: i18n.en, result with true: {"i18n":{"en":{..}}}, result with false: {..}
			"array.on.single.values": true, // result with true: ["single value"], result with false: "single value"
			"depth.level":            0,    // query: i18n, result with 1: {"i18n":{"en":{}}}, result with 0: {"i18n":{"en":{"bla":"en bla"}}}
		},
		"write": {
			"array.replace": false,
		},
	}
	queryWildcardRegex   = regexp.MustCompile(`(\*+)`)
	queryUnifierRegex    = regexp.MustCompile(`\.|\[\d+\]|[^\*\.\[\]\"\']+|[\"\']([^\*\[\]\"\']+)[\"\']`)
	findAllRootKeysRegex = regexp.MustCompile(`^[^\.\[\]]+$`)
	POINT                = []byte(".")
	JSON_OBJECT_START    = []byte("{")
	JSON_ARRAY_START     = []byte("[")
	JSON_ARRAY_END       = []byte("]")
	SINGLE_QUOTE         = []byte("'")
	DOUBLE_QUOTE         = []byte("\"")
)

const (
	JSONArray = 1
	JSONMap   = 2
)

//NiceDB is a fast and a structured database
type NiceDB struct {
	kvDB       *bolt.DB
	kvDBBucket map[string]*bolt.Bucket
	path       string
	options    map[string]map[string]interface{}
}

var niceDBs map[string]*NiceDB

//Open the db on the current path
func Open(path string) (*NiceDB, error) {
	if niceDBs == nil {
		niceDBs = make(map[string]*NiceDB)
	}
	if niceDBs[path] == nil {
		ndb := &NiceDB{}
		niceDBs[path] = ndb
		ndb.options = defaultOptions
		ndb.path = path
		var err error
		ndb.kvDB, err = bolt.Open(ndb.path, 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			return nil, err
		}
		ndb.kvDBBucket = make(map[string]*bolt.Bucket)
	}
	return niceDBs[path], nil
}

type absoluteObj struct {
	objType int
	left    int
	right   int
}

//Read written data
func (ndb *NiceDB) Read(selectors []string, options map[string]interface{}) []interface{} {
	results := make([]interface{}, len(selectors))
	var selector *[]byte
	var keys *[][]byte
	var data interface{}
	absoluteKeyPath := ndb.getOptionValueFor("read", "absolute.key.path", options).(bool)

	// underlying db stuff >>>>>>>
	tx, err := ndb.kvDB.Begin(false)
	if err != nil {
		return results
	}
	region := ndb.getOptionValueFor("global", "region", options).([]byte)
	b := tx.Bucket(region)
	if b == nil {
		return results
	}

	// underlying db stuff end <<<<<

	for i := range selectors {
		if ndb.hasWildcard(selectors[i]) {
			//search
			data = ndb.search(selectors[i], options, b)
		} else {
			//get
			selector, keys = ndb.unify(selectors[i])
			data = ndb._read(selector, b, nil)
			if absoluteKeyPath {
				ndb.makeAbsolute(keys, &data, b)
			}
		}
		results[i] = data
	}

	// underlying db stuff >>>>>>>
	tx.Rollback()
	// underlying db stuff end <<<<<
	return results
}

//TODO optimize makeAbsolute fucntion
func (ndb *NiceDB) makeAbsolute(keys *[][]byte, data *interface{}, bu *bolt.Bucket) {
	pathBuffer := new(bytes.Buffer)
	objectPathBuffer := new(bytes.Buffer)
	var jsonBytes *[]byte
	var getErr error
	var newPath []byte
	l := len(*keys)
	absoluteArray := make([]*absoluteObj, l)
	var aaIndex = 0
	l--
	var k []byte
	var r = 0
	for g := l; g >= 0; g-- {
		k = (*keys)[g]
		if bytes.Equal(k, POINT) {
			continue
		}
		for r = 0; r <= g; r++ {
			k = (*keys)[r]
			pathBuffer.Write(k)

		}
		newPath = pathBuffer.Bytes()
		pathBuffer.Reset()
		jsonBytes, getErr = ndb.get(&newPath, bu)
		if getErr == nil {
			if ndb.isJSONMap(jsonBytes) {
				absoluteArray[aaIndex] = &absoluteObj{objType: JSONMap, left: 0, right: g}
				aaIndex++
			} else if ndb.isJSONArray(jsonBytes) {
				absoluteArray[aaIndex] = &absoluteObj{objType: JSONArray, left: 0, right: g}
				aaIndex++
			}
		}
	}
	l = len(absoluteArray)
	var aa *absoluteObj
	r = 0
	for aaIndex = 0; aaIndex < l; aaIndex++ {
		aa = absoluteArray[aaIndex]
		if aa != nil {
			for r = aaIndex + 1; r < l; r++ {
				if absoluteArray[r] != nil {
					if absoluteArray[r].right-1 >= 0 {
						aa.left = absoluteArray[r].right - 1
						aa.left = ndb.nextLeftKey(keys, &aa.left, &l)
					}
					break
				}
			}
		}
	}
	for aaIndex = 0; aaIndex < l; aaIndex++ {
		aa = absoluteArray[aaIndex]
		if aa != nil {
			switch aa.objType {
			case JSONMap:
				objectPathBuffer.Reset()
				ndb.setNewPath(keys, objectPathBuffer, &aa.left, &aa.right)
				wrapper := make(map[string]interface{})
				s := objectPathBuffer.String()
				wrapper[s] = *data
				*data = wrapper
			case JSONArray:
				objectPathBuffer.Reset()
				ndb.setNewPath(keys, objectPathBuffer, &aa.left, &aa.right)
				wrapper := make([]interface{}, 1)
				wrapper[0] = *data
				*data = wrapper
			}
		}
	}
}

func (ndb *NiceDB) nextLeftKey(keys *[][]byte, from *int, to *int) int {
	for r := *from; r >= 0; r-- {
		if bytes.Equal((*keys)[r], POINT) || bytes.HasPrefix((*keys)[r], JSON_ARRAY_START) {
			continue
		}
		return r
	}
	return 0
}
func (ndb *NiceDB) setNewPath(keys *[][]byte, b *bytes.Buffer, from *int, to *int) {
	for r := *from; r <= *to; r++ {
		b.Write((*keys)[r])
	}
}

//TODO impl wildcard search
func (ndb *NiceDB) search(selector string, options map[string]interface{}, b *bolt.Bucket) []interface{} {
	c := b.Cursor()
	res := make([]interface{}, 0)
	//sel := ndb.wildcardPieces(selector)
	absoluteKeyPath := ndb.getOptionValueFor("read", "absolute.key.path", options).(bool)
	var r interface{}
	var taken [][]byte
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if ndb.alreadyTaken(&k, &taken) {
			continue
		}
		r = ndb._read(&k, b, &v)
		if r != nil {
			taken = append(taken, k)
			if absoluteKeyPath {
				wrapper := make(map[string]interface{})
				wrapper[string(k)] = r
				r = wrapper
			}
			res = append(res, r)
		}
	}
	return res
}

func (ndb *NiceDB) alreadyTaken(k *[]byte, taken *[][]byte) bool {
	for _, item := range *taken {
		if bytes.HasPrefix(*k, item) {
			return true
		}
	}
	return false
}

//Write new or update data
func (ndb *NiceDB) Write(selector map[string]interface{}, options map[string]interface{}) error {
	fmt.Println("before write tx")
	// underlying db stuff >>>>>>>
	tx, err := ndb.kvDB.Begin(true)
	if err != nil {
		return err
	}
	fmt.Println("write tx started")
	defer tx.Rollback()

	b := tx.Bucket(ndb.getOptionValueFor("global", "region", options).([]byte))
	if b == nil {
		b, err = tx.CreateBucket(ndb.getOptionValueFor("global", "region", options).([]byte))
		if err != nil {
			return err
		}
	}
	// underlying db stuff end <<<<<
	var writeKey []byte
	for k := range selector {
		writeKey = []byte(k)
		err = ndb._write(&writeKey, selector[k], options, b)
		if err != nil {
			return err
		}
	}
	fmt.Println("before commit")
	// underlying db stuff >>>>>>>
	if err = tx.Commit(); err != nil {
		return err
	}
	fmt.Println("commited")
	// underlying db stuff end <<<<<
	return nil
}

//Delete specific entries
func (ndb *NiceDB) Delete(selectors []string, options map[string]interface{}) error {
	// underlying db stuff >>>>>>>
	tx, err := ndb.kvDB.Begin(true)
	if err != nil {
		return err
	}
	// underlying db stuff end <<<<<

	b := tx.Bucket(ndb.getOptionValueFor("global", "region", options).([]byte))

	if b != nil {
		for _, sel := range selectors {
			selector, _ := ndb.unify(sel)
			err = ndb._delete(selector, b)
			if err != nil {
				return err
			}
		}

		// underlying db stuff >>>>>>>
		if err = tx.Commit(); err != nil {
			return err
		}
		// underlying db stuff end <<<<<
	}
	tx.Rollback()
	return err
}

func (ndb *NiceDB) _delete(fullKeyPath *[]byte, bu *bolt.Bucket) error {
	b, err := ndb.get(fullKeyPath, bu)
	if err != nil || b == nil {
		return nil
	}
	var data interface{}
	isObjectType, err := ndb.bytesToType(&data, b)
	if !isObjectType {
		return bu.Delete(*fullKeyPath)
	}
	switch v := data.(type) {
	case []interface{}:
		size := int(v[0].(float64))
		if size > 0 {
			var indexBytes []byte
			//var arr = make([]interface{}, size)
			for i := 0; i < size; i++ {
				indexBytes = []byte(strconv.Itoa(i))
				ndb._delete(ndb.concat4(fullKeyPath, &JSON_ARRAY_START, &indexBytes, &JSON_ARRAY_END), bu)
			}
		}
	case map[string]interface{}:
		size := len(v)
		if size > 0 {
			var byteKey []byte
			for mk := range v {
				byteKey = []byte(mk)
				ndb._delete(ndb.concat3(fullKeyPath, &POINT, &byteKey), bu)
			}
		}
	}
	bu.Delete(*fullKeyPath)
	return nil
}

/**
clean path to a unique unit
from
omg.blabla["blablabl"].omg.blabla[1]
to
omg.blabla.blablabl.omg.blabla[1]

from
[23]omg.blabla["myemail@gmail.com"].omg.blabla[1]
to
[23]omg.blabla.myemail@gmail.com.omg.blabla[1]
**/
func (ndb *NiceDB) unify(selector string) (*[]byte, *[][]byte) {
	pieces := queryUnifierRegex.FindAllSubmatch([]byte(selector), -1)
	var cleanPieces [][]byte
	newStrWriter := new(bytes.Buffer)
	for _, match := range pieces {
		if len(match[1]) > 0 {
			cleanPieces = append(cleanPieces, POINT, match[1])
			newStrWriter.Write(POINT)
			newStrWriter.Write(match[1])
		} else {
			cleanPieces = append(cleanPieces, match[0])
			newStrWriter.Write(match[0])
		}
	}
	b := newStrWriter.Bytes()
	return &b, &cleanPieces
}

func (ndb *NiceDB) hasWildcard(sel string) bool {
	return strings.HasPrefix(sel, "*")
}

func (ndb *NiceDB) wildcardPieces(sel string) []byte {
	bts := []byte(sel)
	res := queryWildcardRegex.FindAll(bts, -1)
	fmt.Println(res)
	for i, match := range res {
		fmt.Println(string(match), "found at index", i)
	}
	return []byte(strings.Replace(sel, "*", "", -1))
}

func (ndb *NiceDB) _read(fullKeyPath *[]byte, bu *bolt.Bucket, b *[]byte) interface{} {
	var err error
	if b == nil {
		b, err = ndb.get(fullKeyPath, bu)
	}
	if err != nil || b == nil {
		return nil
	}
	var data interface{}
	isObjectType, err := ndb.bytesToType(&data, b)
	if !isObjectType {
		return data
	}
	//fmt.Printf("_read: %v, %v\n", reflect.TypeOf(data), data)
	switch v := data.(type) {
	case []interface{}:
		size := int(v[0].(float64))
		if size > 0 {
			var arr = make([]interface{}, size)
			var indexBytes []byte
			for i := 0; i < size; i++ {
				indexBytes = []byte(strconv.Itoa(i))
				arr[i] = ndb._read(ndb.concat4(fullKeyPath, &JSON_ARRAY_START, &indexBytes, &JSON_ARRAY_END), bu, nil)
				//ndb.sortAsc(&arr, &i)
			}
			return arr
		}
	case map[string]interface{}:
		size := len(v)
		if size > 0 {
			var byteKey []byte
			for mk := range v {
				byteKey = []byte(mk)
				v[mk] = ndb._read(ndb.concat3(fullKeyPath, &POINT, &byteKey), bu, nil)
			}
			return v
		}
	}
	return data
}

func (ndb *NiceDB) bytesToType(data *interface{}, b *[]byte) (bool, error) {
	var err error
	if ndb.isJSON(b) {
		err = json.Unmarshal(*b, data)
		if err != nil || *data == nil {
			return false, err
		}
	} else {
		dlength := len(*b)
		if dlength == 0 {
			return false, nil
		}
		// covert literal value to go
		if (*b)[0] == '"' { //string
			*data = string(*b)
		} else {
			if dlength > 3 && (*b)[0] == '`' { //literal
				if (*b)[1] == 's' && (*b)[2] == '|' {
					*data = string((*b)[3:])
				} else if (*b)[1] == 'i' { //int
					var bigInt int64
					if (*b)[2] == '|' {
						//int
						bigInt, err = strconv.ParseInt(string((*b)[3:]), 10, 0)
						*data = int(bigInt)
					} else if (*b)[2] == '8' && (*b)[3] == '|' {
						//int8
						bigInt, err = strconv.ParseInt(string((*b)[4:]), 10, 8)
						*data = int8(bigInt)
					} else if (*b)[2] == '1' && (*b)[4] == '|' {
						//int16
						bigInt, err = strconv.ParseInt(string((*b)[5:]), 10, 16)
						*data = int16(bigInt)
					} else if (*b)[2] == '3' && (*b)[4] == '|' {
						//int32
						bigInt, err = strconv.ParseInt(string((*b)[5:]), 10, 32)
						*data = int32(bigInt)
					} else if (*b)[2] == '6' && (*b)[4] == '|' {
						//int64
						bigInt, err = strconv.ParseInt(string((*b)[5:]), 10, 64)
						*data = int64(bigInt)
					}
				} else if (*b)[1] == 'u' && (*b)[2] == 'i' { //uint
					var uBigInt uint64
					if (*b)[3] == '|' {
						//uint
						uBigInt, err = strconv.ParseUint(string((*b)[4:]), 10, 0)
						*data = uint(uBigInt)
					} else if (*b)[3] == '8' && (*b)[4] == '|' {
						//uint8
						uBigInt, err = strconv.ParseUint(string((*b)[5:]), 10, 8)
						*data = uint8(uBigInt)
					} else if (*b)[3] == '1' && (*b)[5] == '|' {
						//uint16
						uBigInt, err = strconv.ParseUint(string((*b)[6:]), 10, 16)
						*data = uint16(uBigInt)
					} else if (*b)[3] == '3' && (*b)[5] == '|' {
						//uint32
						uBigInt, err = strconv.ParseUint(string((*b)[6:]), 10, 32)
						*data = uint32(uBigInt)
					} else if (*b)[3] == '6' && (*b)[5] == '|' {
						//uint64
						uBigInt, err = strconv.ParseUint(string((*b)[6:]), 10, 64)
						*data = uint64(uBigInt)
					}
				} else if (*b)[1] == 'u' && (*b)[2] == 'i' && (*b)[3] == 'p' { //uintptr
					//uintptr
				} else if (*b)[1] == 'f' { //float
					var bigFloat float64
					if (*b)[2] == '3' && (*b)[4] == '|' {
						//float32
						bigFloat, err = strconv.ParseFloat(string((*b)[5:]), 32)
						*data = float32(bigFloat)
					} else if (*b)[2] == '6' && (*b)[4] == '|' {
						//float64
						bigFloat, err = strconv.ParseFloat(string((*b)[5:]), 64)
						*data = float64(bigFloat)
					}
				} else if (*b)[1] == 'c' { //complex
					var bigFloat float64
					if (*b)[2] == '6' && (*b)[4] == '|' {
						//complex64
						bigFloat, err = strconv.ParseFloat(string((*b)[5:]), 32)
						*data = float32(bigFloat)
					} else if (*b)[2] == '1' && (*b)[5] == '|' {
						//complex128
						bigFloat, err = strconv.ParseFloat(string((*b)[5:]), 64)
						*data = float64(bigFloat)
					}
				} else if (*b)[1] == 'b' && (*b)[2] == '|' { //bool
					if (*b)[2] == 't' {
						*data = true
					} else {
						*data = false
					}
				}
			}
			if err != nil {
				*data = nil
			}
		}
		//fmt.Printf("_read: %v, %v\n", reflect.TypeOf(data), data)
		return false, nil
	}
	return true, nil
}

//func (ndb *NiceDB) sortAsc(a *[]interface{}, index *int){
//	if((*index) > 0){
//		fmt.Println("sortAsc", (*a)[(*index)-1] < (*a)[(*index)])
//	}
//	//switch v := (*a)[(*index)].(type) {
//	//case []interface{}:
//	//case map[string]interface{}:
//	//}
//}

func (ndb *NiceDB) concat3(a *[]byte, b *[]byte, c *[]byte) *[]byte {
	newPath := new(bytes.Buffer)
	newPath.Write(*a)
	newPath.Write(*b)
	newPath.Write(*c)
	by := newPath.Bytes()
	return &by
}

func (ndb *NiceDB) concat4(a *[]byte, b *[]byte, c *[]byte, d *[]byte) *[]byte {
	newPath := new(bytes.Buffer)
	newPath.Write(*a)
	newPath.Write(*b)
	newPath.Write(*c)
	newPath.Write(*d)
	by := newPath.Bytes()
	return &by
}

func (ndb *NiceDB) _write(fullKeyPath *[]byte, data interface{}, options map[string]interface{}, b *bolt.Bucket) error {
	v := reflect.ValueOf(data)
	var value []byte
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			value = []byte("`b|t")
		} else {
			value = []byte("`b|f")
		}
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Int:
		value = []byte(fmt.Sprintf("`i|%v", data))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Int8:
		value = []byte(fmt.Sprintf("`i8|%v", data))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Int16:
		value = []byte(fmt.Sprintf("`i16|%v", data))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Int32:
		value = []byte(fmt.Sprintf("`i32|%v", data))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Int64:
		value = []byte(fmt.Sprintf("`i64|%v", data))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Uint:
		value = []byte(fmt.Sprintf("`ui|%v", data))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Uint8:
		value = []byte(fmt.Sprintf("`ui8|%v", data))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Uint16:
		value = []byte(fmt.Sprintf("`ui16|%v", data))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Uint32:
		value = []byte(fmt.Sprintf("`ui32|%v", data))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Uint64:
		value = []byte(fmt.Sprintf("`ui64|%v", data))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Uintptr:
		fmt.Printf("`uip|%v", data)
		value = []byte(fmt.Sprintf("`uip|%v", data))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Float32:
		value = []byte("`f32|" + strings.TrimSpace(fmt.Sprintf("%30.30G", data)))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Float64:
		value = []byte("`f64|" + strings.TrimSpace(fmt.Sprintf("%30.30G", data)))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.String:
		value = []byte(fmt.Sprintf("`s|%v", v))
		return ndb.put(fullKeyPath, &value, b)
	case reflect.Array, reflect.Slice:
		if v.IsNil() {
			return nil
		}
		var newLen int
		newLen = v.Len()
		if newLen == 0 {
			return nil
		}
		var startIndex, oldLen int
		var err error
		if !ndb.getOptionValueFor("write", "array.replace", options).(bool) {
			var d *[]byte
			d, err = ndb.get(fullKeyPath, b)
			if err == nil && d != nil {
				var f []int
				err = json.Unmarshal(*d, &f)
				if err == nil && f != nil && len(f) == 1 {
					oldLen = f[0]
					newLen = oldLen + newLen
					startIndex = oldLen
				}
			}
		}
		value = []byte(strconv.Itoa(newLen))
		err = ndb.put(fullKeyPath, ndb.concat3(&JSON_ARRAY_START, &value, &JSON_ARRAY_END), b)
		if err != nil {
			return err
		}
		var indexBytes []byte
		for i := 0; i < v.Len(); i++ {
			indexBytes = []byte(strconv.Itoa(startIndex))
			err = ndb._write(ndb.concat4(fullKeyPath, &JSON_ARRAY_START, &indexBytes, &JSON_ARRAY_END), v.Index(i).Interface(), options, b)
			startIndex++
			if err != nil {
				return err
			}
		}
		return err
	case reflect.Map:
		if v.IsNil() {
			return nil
		}
		if v.Len() == 0 {
			return nil
		}
		d, err := ndb.get(fullKeyPath, b)
		var obj map[string]byte
		if err == nil && d != nil && len(*d) > 0 {
			var f interface{}
			err = json.Unmarshal(*d, &f)
			if err != nil {
				return err
			}
			if tmpObj, ok := f.(map[string]byte); ok {
				obj = tmpObj
			} else {
				obj = make(map[string]byte)
			}
		} else {
			obj = make(map[string]byte)
		}

		//set all keys available in the given map
		var byteKey []byte

		for _, mk := range v.MapKeys() {
			obj[mk.String()] = 1
			byteKey = []byte(mk.String())
			ndb._write(ndb.concat3(fullKeyPath, &POINT, &byteKey), v.MapIndex(mk).Interface(), options, b)
		}
		var marshaledBytes []byte
		marshaledBytes, err = json.Marshal(obj)
		if err != nil {
			return err
		}
		return ndb.put(fullKeyPath, &marshaledBytes, b)
	default:
		fmt.Println("Kind: ", v.Kind(), v, data)
		fmt.Println("I don't know, ask stackoverflow.")
	}
	return nil
}

func (ndb *NiceDB) put(key *[]byte, value *[]byte, b *bolt.Bucket) error {
	if err := b.Put(*key, *value); err != nil {
		return err
	}
	return nil
}

func (ndb *NiceDB) get(key *[]byte, b *bolt.Bucket) (*[]byte, error) {
	myBytes := b.Get(*key)
	return &myBytes, nil
}

//Dump to show whats in the db
func (ndb *NiceDB) Dump(options map[string]interface{}) {
	tx, err := ndb.kvDB.Begin(false)
	if err != nil {
		return
	}
	defer tx.Rollback()
	b := tx.Bucket(ndb.getOptionValueFor("global", "region", options).([]byte))
	if b == nil {
		fmt.Println("empty")
		return
	}
	c := b.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		fmt.Printf("%s=%s\n", k, v)
	}
}

func (ndb *NiceDB) getOptionValueFor(category string, key string, options map[string]interface{}) interface{} {
	if options != nil {
		if val, ok := options[key]; ok {
			return val
		}
	}
	return ndb.options[category][key]
}

func (ndb *NiceDB) getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (ndb *NiceDB) isJSON(json *[]byte) bool {
	return ndb.isJSONMap(json) || ndb.isJSONArray(json)
}

func (ndb *NiceDB) isJSONMap(json *[]byte) bool {
	return json != nil && bytes.HasPrefix(*json, JSON_OBJECT_START)
}

func (ndb *NiceDB) isJSONArray(json *[]byte) bool {
	return json != nil && bytes.HasPrefix(*json, JSON_ARRAY_START)
}

func (ndb *NiceDB) isGoMap(data interface{}) bool {
	_, ok := data.(map[string]interface{})
	return ok
}

func (ndb *NiceDB) isGoArray(data interface{}) bool {
	_, ok := data.([]interface{})
	return ok
}

func (ndb *NiceDB) isQuoted(data *[]byte) bool {
	return bytes.HasPrefix(*data, SINGLE_QUOTE) || bytes.HasPrefix(*data, DOUBLE_QUOTE)
}

func (ndb *NiceDB) UUID() (uuid string) {
	n := time.Now()
	b := make([]byte, 18)
	_, err := rand.Read(b)
	if err != nil {
		uuid = fmt.Sprintf("%X%X%v%v%X", n.Nanosecond()*10, n.Unix(), n.Nanosecond(), n.Unix(), time.Now().UnixNano())
		return
	}
	uuid = fmt.Sprintf("%X%X%X%X%X%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:], n.UnixNano())
	return
}

//Close db
func (ndb *NiceDB) Close() {
	ndb.kvDB.Close()
}
