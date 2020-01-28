package sys

import (
	"io"
	"path/filepath"

	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

type Session struct {
	S     *model.Session
	db    *storage.DBSet
	cache *cache.Cache
}

func NewSession(sys *System, usr *model.User) (*Session, error) {
	s := &model.Session{
		ID:         uuid.NewV4().String(),
		Rights:     usr.Role,
		UsrID:      usr.ID,
		UserName:   usr.Name,
		SessionDir: usr.ID + "_" + uuid.NewV4().String(),
	}
	err := sys.DB.Session.Put(s)
	return &Session{S: s, db: sys.DB, cache: sys.cache}, err
}

var _ model.Auth = (*Session)(nil)

func (s *Session) UserID() string {
	return s.S.UsrID
}

func (s *Session) AccessRights() model.Role {
	if s.S.UserName == "anonymous" {
		return model.USER
	}
	return s.S.Rights
}

func (s *Session) GetSessionDir() string {
	return s.S.SessionDir
}

func (s *Session) FilePath(filename string) string {
	return filepath.Join(s.S.SessionDir, filename)
}

func (s *Session) WriteFile(filename string, reader io.Reader) error {
	return s.db.Files.Write(s.FilePath(filename), reader)
}

func (s *Session) DeleteFile(filename string) error {
	return s.db.Files.Delete(s.FilePath(filename))
}

func (s *Session) Delete(k string) {
	s.db.Session.DeleteValue(s.key(k))
}

func (s *Session) Put(k string, val interface{}) {
	s.db.Session.PutValue(s.key(k), val)
}

func (s *Session) Get(k string, ref interface{}) error {
	return s.db.Session.GetValue(s.key(k), ref)
}

func (s *Session) GetMemory(k string) (interface{}, bool) {
	k2 := s.key(k)
	v, ok := s.cache.Get(k2)
	if ok {
		// touch to prolong expiration
		s.cache.SetDefault(k2, v)
	}
	return v, ok
}

func (s *Session) PutMemory(k string, val interface{}) {
	s.cache.SetDefault(s.key(k), val)
}

func (s *Session) DeleteMemory(k string) {
	s.cache.Delete(s.key(k))
}

func (s *Session) DeleteAll() error {
	return s.db.Session.Delete(s.S)
}

func (s *Session) key(k string) string {
	return s.S.ID + "_" + s.S.UsrID + "_" + k
}
