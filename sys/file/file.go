/*
This package simplifies file IO and handles the file path when using JSON or MSGPACK.
When using Output the path will be the URI like /api/my/file/uri/{unique_id}
When using Input the path will be the relative file system path my/relative/file/path/{unique_id}
*/
package file

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"

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

func FromJSONBytes(baseFilePath string, bts []byte) (*IO, error) {
	f := &IO{baseFileDir: baseFilePath}
	err := f.UnmarshalJSON(bts)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func FromMap(baseFilePath string, m map[string]interface{}) *IO {
	f := &IO{baseFileDir: baseFilePath}
	fromMap(f, m)
	return f
}

func New(baseDir string, meta Meta) *IO {
	return &IO{
		baseFileDir: baseDir,
		path:        RndUUID(),
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

func (me *IO) Ref() string {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return me.ref
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

func (me *IO) Stat() (os.FileInfo, error) {
	p := me.Path()
	if p != "" {
		return os.Stat(p)
	}
	return nil, os.ErrNotExist
}

func (me *IO) ToMap(pathNameOnly bool) map[string]interface{} {
	me.lock.RLock()
	defer me.lock.RUnlock()
	if pathNameOnly {
		return map[string]interface{}{
			"ref":         me.ref,
			"contentType": me.meta.ContentType,
			"name":        me.meta.Name,
			"size":        me.meta.Size,
			"path":        me.path,
			"hash":        me.Hash,
		}
	}
	return map[string]interface{}{
		"ref":         me.ref,
		"contentType": me.meta.ContentType,
		"name":        me.meta.Name,
		"size":        me.meta.Size,
		"path":        filepath.Join(me.baseFileDir, me.path),
		"hash":        me.Hash,
	}
}
func (me *IO) ContentType() string {
	me.lock.RLock()
	defer me.lock.RUnlock()
	return me.meta.ContentType
}

func (me *IO) Move(oldPath string) error {
	me.lock.Lock()
	defer me.lock.Unlock()
	p := me.Path()
	if p == "" || oldPath == "" {
		return os.ErrInvalid
	}
	f, err := os.Stat(oldPath)
	if err == nil {
		me.meta.Size = f.Size()
	}
	return os.Rename(oldPath, p)
}

func (me *IO) Write(reader io.Reader) (int64, error) {
	me.lock.Lock()
	defer me.lock.Unlock()
	var f *os.File
	f, err := os.OpenFile(me.Path(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		return 0, err
	}
	defer f.Close()
	me.meta.Size, err = io.Copy(f, reader)
	return me.meta.Size, err
}

func (me *IO) CpTo(to *IO) (int64, error) {
	var f *os.File
	var f2 *os.File
	me.lock.RLock()
	defer func() {
		if f != nil {
			f.Close()
		}
		if f2 != nil {
			f2.Close()
		}
		me.lock.RUnlock()
	}()

	f, err := os.OpenFile(me.Path(), os.O_RDONLY, 0600)
	if err != nil {
		return 0, err
	}
	var fstat os.FileInfo
	fstat, err = f.Stat()
	if err != nil {
		return 0, err
	}
	me.meta.Size = fstat.Size()

	f2, err = os.OpenFile(to.Path(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if os.IsExist(err) {
		err = nil
	}
	return io.CopyN(f2, f, me.meta.Size)

}

func (me *IO) Read(writer io.Writer) (int64, error) {
	var f *os.File
	me.lock.RLock()
	defer func() {
		if f != nil {
			f.Close()
		}
		me.lock.RUnlock()
	}()

	f, err := os.OpenFile(me.Path(), os.O_RDONLY, 0600)
	if err != nil {
		return 0, err
	}
	var fstat os.FileInfo
	fstat, err = f.Stat()
	if err != nil {
		return 0, err
	}
	me.meta.Size = fstat.Size()
	return io.CopyN(writer, f, me.meta.Size)
}

func (me *IO) ReadAll() ([]byte, error) {
	bt := &bytes.Buffer{}
	_, err := me.Read(bt)
	if err != nil {
		return nil, err
	}
	return bt.Bytes(), nil
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

func RndUUID() string {
	u2 := uuid.NewV4()
	return u2.String()
}
