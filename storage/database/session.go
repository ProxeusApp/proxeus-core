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

// NewSessionDB returns a handle to the session & token database with the supplied expiration times
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

// Get returns a session
func (d *SessionDB) Get(sid string) (*model.Session, error) {
	var s model.Session
	err := d.db.Get(sessionBucket, sid, &s)
	return &s, err
}

// Put inserts a session
func (d *SessionDB) Put(s *model.Session) error {
	return d.db.Set(sessionBucket, s.ID, s, db.OptionWithTTL(d.sessionExpiration))
}

// Delete removes a session
func (d *SessionDB) Delete(s *model.Session) error {
	return d.db.Delete(sessionBucket, s.ID)
}

// GetTokenRequest returns a token request
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

// PutTokenRequest inserts a token request
func (d *SessionDB) PutTokenRequest(r *model.TokenRequest) error {
	return d.db.Set(tokenBucket, r.Token, r, db.OptionWithTTL(d.tokenExpiration))
}

// DeleteTokenRequest removes a token request
func (d *SessionDB) DeleteTokenRequest(r *model.TokenRequest) error {
	return d.db.Delete(tokenBucket, r.Token)
}

// GetValue returns the value for the provided key
func (d *SessionDB) GetValue(key string, v interface{}) error {
	return d.db.Get(valueBucket, key, v)
}

// PutValue sets the value for a key
func (d *SessionDB) PutValue(key string, v interface{}) error {
	return d.db.Set(valueBucket, key, v, db.OptionWithTTL(d.sessionExpiration))
}

// DeleteValue removes the value of a key
func (d *SessionDB) DeleteValue(key string) error {
	return d.db.Delete(valueBucket, key)
}

func (d *SessionDB) Close() error { return d.db.Close() }
