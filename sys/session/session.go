package session

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ProxeusApp/proxeus-core/sys/cache"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type (
	Session struct {
		id         string        //read only
		rights     model.Role    //read only
		userID     string        //read only
		userName   string        //read only
		sessionDir string        //read only
		store      *cache.UCache //managed read and write storage
		manager    *Manager
	}
)

func (me *Session) init(manager *Manager) {
	me.manager = manager
	config := cache.UCacheConfig{
		ValueBehaviour: &cache.ValueBehaviour{
			ValueBeType:        cache.CallOnLoadAndOnExpire,
			OnLoadMethodName:   "OnLoad",
			OnExpireMethodName: "Close",
		},
		DiskStorePath: filepath.Join(me.SessionDir(), "cache"),
		StoreType:     cache.MemAndDiskCache,
		NoExpiry:      true,
	}
	var err error
	me.store, err = cache.NewUCache(config)
	if err != nil {
		log.Println("error with MemAndDiskCache", err)
		config.StoreType = cache.MemCache
		log.Println("trying MemCache")
		me.store, err = cache.NewUCache(config)
		if err != nil {
			log.Println("error with MemCache too", err)
		}
	}
}

//loaded is called by the manager when the session gets loaded from the disk
func (me *Session) loaded(manager *Manager) {
	log.Println("loaded", me.id)
	me.init(manager)
}

//ID of this session
func (me *Session) ID() string {
	return me.id
}

//AccessRights provides the users permission role and implements the Auth interface
func (me *Session) AccessRights() model.Role {
	return me.rights
}

//UserID provides a unique user identification and implements the Auth interface
func (me *Session) UserID() string {
	return me.userID
}

func (me *Session) SetUserID(userID string) {
	me.userID = userID
}

//UserName provides a readable user identification
func (me *Session) UserName() string {
	return me.userName
}

//SessionDir provides the sessions directory
func (me *Session) SessionDir() string {
	me.ensureSessionDir()
	return me.sessionDir
}

//ReadFile from the sessions directory with filename and writer
func (me *Session) ReadFile(filename string, writer io.Writer) (n int64, err error) {
	var f *os.File
	f, err = os.OpenFile(filepath.Join(me.sessionDir, filename), os.O_RDONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	var fstat os.FileInfo
	fstat, err = f.Stat()
	if err != nil {
		return
	}
	n, err = io.CopyN(writer, f, fstat.Size())
	return
}

//FilePath concats the sessions directory with filename
func (me *Session) FilePath(filename string) (fpath string, exists bool) {
	fpath = filepath.Join(me.sessionDir, filename)
	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		exists = false
		return
	}
	exists = true
	return
}

//WriteFile to the sessions directory with filename and reader
func (me *Session) WriteFile(filename string, reader io.Reader) (written int64, err error) {
	tmpFilename := filename + "_" + strconv.Itoa(time.Now().Nanosecond())
	tmpPath := filepath.Join(me.sessionDir, tmpFilename)
	defer func() {
		//cleanup if not successfully written and moved
		os.Remove(tmpPath)
	}()
	err = me.ensureSessionDir()
	if err != nil {
		return
	}
	var f *os.File
	f, err = os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	written, err = io.Copy(f, reader)
	f.Close()
	if err != nil {
		return
	}
	err = me.MoveFile(tmpFilename, filepath.Join(me.sessionDir, filename))
	return
}

func (me *Session) MoveFile(filename, newFullPath string) error {
	return os.Rename(filepath.Join(me.sessionDir, filename), newFullPath)
}

//DeleteFile from sessions directory with filename
func (me *Session) DeleteFile(filename string) error {
	return os.RemoveAll(filepath.Join(me.sessionDir, filename))
}

func (me *Session) ensureSessionDir() error {
	if me.manager != nil {
		return me.manager.ensureDir(me.sessionDir)
	}
	return fmt.Errorf("couldn't ensure session dir exists, manager reference missing")
}

func (me *Session) String() string {
	return fmt.Sprintf(`{"id":"%s","rights":%d,"userId":"%s","userName":"%s","sessionDir":"%s"}`, me.id, me.rights, me.userID, me.userName, me.sessionDir)
}

func (me *Session) MarshalJSON() ([]byte, error) {
	return []byte(me.String()), nil
}

func (me *Session) UnmarshalJSON(b []byte) error {
	var jsonObj struct {
		ID         string     `json:"id"`
		Rights     model.Role `json:"rights"`
		UserID     string     `json:"userId"`
		UserName   string     `json:"userName"`
		SessionDir string     `json:"sessionDir"`
	}
	err := json.Unmarshal(b, &jsonObj)
	if err != nil {
		return err
	}
	me.id = jsonObj.ID
	me.rights = jsonObj.Rights
	me.userID = jsonObj.UserID
	me.userName = jsonObj.UserName
	me.sessionDir = jsonObj.SessionDir
	return nil
}

//Delete and Close data on memory
func (me *Session) Delete(key string) {
	me.store.Close()
	if me.store != nil {
		me.store.Remove(key)
	}
}

//Put data for memory only
func (me *Session) Put(key string, val interface{}) {
	if me.store != nil {
		me.store.Put(key, val)
	}
}

//Get data from memory only
func (me *Session) Get(key string, ref interface{}) error {
	if me.store != nil {
		return me.store.Get(key, ref)
	}
	return cache.ErrNotFound
}

func (me *Session) expire() {
	if me.store != nil {
		me.store.Clean()
		me.store.Close()
	}
}

func (me *Session) Kill() error {
	me.expire()
	if me.manager != nil {
		return me.manager.Remove(me.id)
	}
	return fmt.Errorf("couldn't remove session, manager reference missing")
}

//Close removes the session from memory and disk
func (me *Session) close() (err error) {
	if me.store != nil {
		me.store.IterateMemStorage(func(key string, val interface{}) {
			//ensure all changes during runtime are persisted before closing
			me.store.UpdatedValueRef(key)
			if closer, ok := val.(io.Closer); ok {
				closer.Close()
			}
		})
		me.store.Close()
	}
	return nil
}
