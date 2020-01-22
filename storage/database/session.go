package database

import (
	"errors"
	"path"
	"path/filepath"
	"time"

	"github.com/ProxeusApp/proxeus-core/sys/model"

	"github.com/ProxeusApp/proxeus-core/storage/database/db"
)

type SessionDB struct {
	db                db.DB
	sessionExpiration time.Duration
	tokenExpiration   time.Duration
}

const sessionBucket = "sessions"
const tokenBucket = "tokens"
const valueBucket = "values"

func NewSessionDB(c DBConfig, sessionExpiration, tokenExpiration time.Duration) (*SessionDB, error) {
	baseDir := path.Join(c.Dir, "session")
	db, err := db.OpenDatabase(c.Engine, c.URI, filepath.Join(baseDir, "sessions"))
	if err != nil {
		return nil, err
	}
	return &SessionDB{
		db:                db,
		sessionExpiration: sessionExpiration,
		tokenExpiration:   tokenExpiration,
	}, nil
}

func (d *SessionDB) Get(sid string) (*model.Session, error) {
	var s model.Session
	err := d.db.Get(sessionBucket, sid, &s)
	return &s, err
}

func (d *SessionDB) Put(s *model.Session) error {
	return d.db.Set(sessionBucket, s.ID, s, db.OptionWithTTL(d.sessionExpiration))
}

func (d *SessionDB) Delete(s *model.Session) error {
	return d.db.Delete(sessionBucket, s.ID)
}

func (d *SessionDB) GetTokenRequest(t model.TokenType, id string) (*model.TokenRequest, error) {
	var s model.TokenRequest
	err := d.db.Get(tokenBucket, id, &s)
	if err != nil {
		return nil, err
	}
	if s.Type != t {
		return nil, errors.New("token type mismatch")
	}
	return &s, nil
}
func (d *SessionDB) PutTokenRequest(r *model.TokenRequest) error {
	return d.db.Set(tokenBucket, r.Token, r, db.OptionWithTTL(d.tokenExpiration))
}

func (d *SessionDB) DeleteTokenRequest(r *model.TokenRequest) error {
	return d.db.Delete(tokenBucket, r.Token)
}

func (d *SessionDB) GetValue(key string, v interface{}) error {
	return d.db.Get(valueBucket, key, v)
}
func (d *SessionDB) PutValue(key string, v interface{}) error {
	return d.db.Set(valueBucket, key, v, db.OptionWithTTL(d.sessionExpiration))
}
func (d *SessionDB) DeleteValue(key string) error {
	return d.db.Delete(valueBucket, key)
}

func (d *SessionDB) Close() error { return d.db.Close() }
