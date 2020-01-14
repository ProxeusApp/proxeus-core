/*
This package simplifies file IO and handles the file path when using JSON or MSGPACK.
When using Output the path will be the URI like /api/my/file/uri/{unique_id}
When using Input the path will be the relative file system path my/relative/file/path/{unique_id}
*/
package file

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/asdine/storm/codec/msgpack"
	uuid "github.com/satori/go.uuid"
)

var fileExtReplacer = regexp.MustCompile(`(\.\w{2,4}$|$)`)

type (
	Meta struct {
		Name        string
		ContentType string
		Size        int64
		Hash        string
	}
	diskIO struct {
		Ref         string `json:"ref,omitempty" msgpack:"ref,omitempty"`
		Path        string `json:"path" msgpack:"path"`
		Name        string `json:"name" msgpack:"name"`
		ContentType string `json:"contentType" msgpack:"contentType"`
		Size        int64  `json:"size" msgpack:"size"`
		Hash        string `json:"hash"`
	}
	IO struct {
		ref         string
		path        string
		meta        Meta
		baseFileDir string
		lock        sync.RWMutex
		Hash        string
	}
)

var int64Type = reflect.TypeOf(int64(0))

func FromMap(baseFilePath string, m map[string]interface{}) *IO {
	f := &IO{baseFileDir: baseFilePath}
	fromMap(f, m)
	return f
}

func New(baseDir string, meta Meta) *IO {
	return &IO{
		baseFileDir: baseDir,
		path:        uuid.NewV4().String(),
		meta:        meta,
	}
}

func IsFileInfo(maybeFile map[string]interface{}) bool {
	if maybeFile != nil {
		if _, ok := maybeFile["path"]; ok {
			if _, ok := maybeFile["name"]; ok {
				if _, ok := maybeFile["size"]; ok {
					if _, ok := maybeFile["contentType"]; ok {
						return true
					}
				}
			}
		}
	}
	return false
}

func (me *IO) Name() string {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return me.meta.Name
}

func (me *IO) SetRef(ref string) {
	me.lock.RLock()
	defer me.lock.RUnlock()
	me.ref = ref
}

func (me *IO) SetBaseDir(baseDir string) {
	me.lock.RLock()
	defer me.lock.RUnlock()
	me.baseFileDir = baseDir
}

func (me *IO) Meta() Meta {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return me.meta
}

func (me *IO) Update(name, contentType string) {
	me.lock.Lock()
	me.meta.Name = name
	me.meta.ContentType = contentType
	me.lock.Unlock()
}

func (me *IO) Path() string {
	if me.baseFileDir != "" {
		return filepath.Join(me.baseFileDir, me.path)
	}
	return me.path
}

func (me *IO) PathName() string {
	return me.path
}

func (me *IO) Size() int64 {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return me.meta.Size
}

func (me *IO) SetSize(s int64) {
	me.lock.Lock()
	defer me.lock.Unlock()
	me.meta.Size = s
}

func (me *IO) ToMap() map[string]interface{} {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return map[string]interface{}{
		"ref":         me.ref,
		"contentType": me.meta.ContentType,
		"name":        me.meta.Name,
		"size":        me.meta.Size,
		"path":        me.path,
		"hash":        me.Hash,
	}
}
func (me *IO) ContentType() string {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return me.meta.ContentType
}

func (me *IO) NameWithExt(ext string) string {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return NameWithExt(me.meta.Name, ext)
}

func NameWithExt(fileName string, ext string) string {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	return fileExtReplacer.ReplaceAllString(fileName, ext)
}

func (me *IO) String() string {
	bts, err := me.MarshalJSON()
	if err != nil {
		return ""
	}
	return string(bts)
}

func (me *IO) MarshalJSON() ([]byte, error) {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return json.Marshal(diskIO{
		Ref:         me.ref,
		Path:        me.path,
		Name:        me.meta.Name,
		ContentType: me.meta.ContentType,
		Size:        me.meta.Size,
		Hash:        me.Hash,
	})
}

func (me *IO) MarshalMsgpack() ([]byte, error) {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return msgpack.Codec.Marshal(diskIO{
		Ref:         me.ref,
		Path:        me.path,
		Name:        me.meta.Name,
		ContentType: me.meta.ContentType,
		Size:        me.meta.Size,
		Hash:        me.Hash,
	})
}

func (me *IO) UnmarshalMsgpack(b []byte) error {
	var obj *diskIO
	err := msgpack.Codec.Unmarshal(b, &obj)
	if err != nil {
		return err
	}
	fromDiskIO(me, obj)
	return nil
}

func (me *IO) UnmarshalJSON(b []byte) error {
	var jsonObj *diskIO
	err := json.Unmarshal(b, &jsonObj)
	if err != nil {
		return err
	}
	fromDiskIO(me, jsonObj)
	return nil
}

func (me *IO) MarshalBSON() ([]byte, error) {
	d := diskIO{
		Ref:         me.ref,
		Path:        me.path,
		Name:        me.meta.Name,
		ContentType: me.meta.ContentType,
		Size:        me.meta.Size,
		Hash:        me.Hash,
	}
	return bson.Marshal(&d)
}

func (me *IO) UnmarshalBSON(raw []byte) error {
	obj := new(diskIO)
	err := bson.Unmarshal(raw, &obj)
	if err != nil {
		return err
	}
	fromDiskIO(me, obj)
	return nil
}

func fromDiskIO(me *IO, obj *diskIO) {
	if obj.Path != "" {
		me.path = filepath.Base(obj.Path)
	}
	me.ref = obj.Ref
	me.meta.Name = obj.Name
	me.meta.ContentType = obj.ContentType
	me.meta.Size = obj.Size
}

func fromMap(me *IO, jsonObj map[string]interface{}) {
	if len(jsonObj) > 0 {
		if iname, ok := jsonObj["path"]; ok {
			me.path, ok = iname.(string)
			if me.path != "" {
				me.path = filepath.Base(me.path)
			}
		}
		if iname, ok := jsonObj["ref"]; ok {
			me.ref = iname.(string)
		}
		if iname, ok := jsonObj["name"]; ok {
			me.meta.Name, ok = iname.(string)
		}
		if icontentType, ok := jsonObj["contentType"]; ok {
			me.meta.ContentType, ok = icontentType.(string)
		}
		if s, ok := jsonObj["size"]; ok {
			v := reflect.ValueOf(s)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
				v = v.Convert(int64Type)
				me.meta.Size = v.Int()
			}
		}
	}
}

type InMemoryFileInfo struct {
	Path string
	Len  int
}

func (fi InMemoryFileInfo) Name() string       { return fi.Path }
func (fi InMemoryFileInfo) Size() int64        { return int64(fi.Len) }
func (fi InMemoryFileInfo) Mode() os.FileMode  { return 0777 }
func (fi InMemoryFileInfo) ModTime() time.Time { return time.Now() }
func (fi InMemoryFileInfo) IsDir() bool        { return false }
func (fi InMemoryFileInfo) Sys() interface{}   { return nil }
