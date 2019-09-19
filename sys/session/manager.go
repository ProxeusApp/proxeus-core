package session

import (
	"os"
	"path/filepath"
	"time"

	"log"

	uuid "github.com/satori/go.uuid"

	"git.proxeus.com/core/central/sys/cache"
	"git.proxeus.com/core/central/sys/model"
)

type (
	//Notify can be used for statistic or monitoring purposes
	Notify interface {
		//OnSessionCreated is called when manager.New is creating a new session
		OnSessionCreated(id string, s *Session)
		//OnSessionLoaded is called when the session gets loaded from the disk
		OnSessionLoaded(id string, s *Session)
		//OnSessionExpired is called when the session expires. The expiry can be provided when creating the manager.
		//Heads up, session can be nil if it wasn't in use during expiry.
		OnSessionExpired(id string, s *Session)
		//OnSessionRemoved is called when session.Close() or manager.Remove(id) is called.
		OnSessionRemoved(id string)
	}
	Manager struct {
		sessionsDB  *cache.UCache
		sessionsDir string
		notify      Notify
	}
)

const sessions = "sessions"

//NewManager creates a new Manager without a Notify object.
//sessionDir must be a directory, if it is empty the a new directory will be created under the os temp dir.
//expiry must be provided and it must be higher than 0
func NewManager(sessionsDir string, expiry time.Duration) (sm *Manager, err error) {
	return NewManagerWithNotify(sessionsDir, expiry, nil)
}

//NewManagerWithNotify is the same as NewManager but it provides the possibility to pass a Notify object.
func NewManagerWithNotify(sessionsDir string, expiry time.Duration, notify Notify) (sm *Manager, err error) {
	if sessionsDir == "" {
		sessionsDir = os.TempDir()
	}
	sessionsDir = filepath.Join(sessionsDir, sessions)
	s := &Manager{sessionsDir: sessionsDir, notify: notify}
	err = s.ensureDir(sessionsDir)
	if err != nil {
		return
	}
	config := cache.UCacheConfig{
		ValueBehaviour: &cache.ValueBehaviour{
			ValueBeType:  cache.CallOnLoadAndOnExpire,
			OnLoadFunc:   s.onSessionLoad,
			OnExpireFunc: s.onSessionExpire,
		},
		DiskStorePath: filepath.Join(sessionsDir, sessions),
		StoreType:     cache.MemAndDiskCacheSpeedMode,
		ExtendExpiry:  true,
		DefaultExpiry: expiry,
	}
	s.sessionsDB, err = cache.NewUCache(config)
	if err != nil {
		return
	}
	sm = s
	return
}

//New creates a new session, it puts the session into the persistent storage and calls notify.OnSessionCreated if the notify object is provided.
func (me *Manager) New(userId string, usrName string, rights model.Role) (session *Session, err error) {
	sessUUID := uuid.NewV4()
	session = &Session{}
	session.id = sessUUID.String()
	session.userName = usrName
	session.userID = userId
	session.rights = rights
	err = me.setup(session)
	session.init(me)
	if err != nil {
		return nil, err
	}
	err = me.sessionsDB.Put(session.id, session)
	if err != nil {
		return nil, err
	}
	if me.notify != nil {
		me.notify.OnSessionCreated(session.id, session)
	}
	return session, err
}

//SessionsDir provides the session managers main directory for all sessions.
//This method should be used for read only purposes otherwise it can cause damage.
func (me *Manager) SessionsDir() string {
	return me.sessionsDir
}

//Get a session by id either directly from memory or load it from the disk.
//Loading it from disk will trigger notify.OnSessionLoaded
func (me *Manager) Get(id string) (session *Session, err error) {
	err = me.sessionsDB.Get(id, &session)
	return
}

//Remove a session by id
func (me *Manager) Remove(id string) error {
	err := me.sessionsDB.Remove(id)
	if err != nil {
		return err
	}
	if me.notify != nil {
		me.notify.OnSessionRemoved(id)
	}
	return me.remSessionDir(id)
}

func (me *Manager) ensureDir(sessionsDir string) error {
	_, err := os.Stat(sessionsDir)
	if os.IsNotExist(err) {
		perm := os.FileMode(0750)
		err = os.MkdirAll(sessionsDir, perm)
		if err == nil {
			//Mkdir err not important as we have the MkdirAll error already
			err = os.Mkdir(sessionsDir, perm)
			if os.IsExist(err) {
				err = nil
			}
		}
		return err
	}
	return nil
}

func (me *Manager) onSessionLoad(key string, val interface{}) {
	if sess, ok := val.(*Session); ok {
		sess.loaded(me)
		if me.notify != nil {
			me.notify.OnSessionLoaded(key, sess)
		}
	}
}

func (me *Manager) onSessionExpire(key string, val interface{}) {
	if sess, ok := val.(*Session); ok && sess != nil {
		if me.notify != nil {
			me.notify.OnSessionExpired(key, sess)
		}
		sess.expire()
	} else {
		if me.notify != nil {
			me.notify.OnSessionExpired(key, nil)
		}
	}
	_ = me.remSessionDir(key)
}

func (me *Manager) remSessionDir(id string) error {
	err := os.RemoveAll(me.constructSessionDir(id))
	if err != nil {
		log.Println("Session Manager error when removing session directory:", err)
	}
	return err
}

func (me *Manager) setup(session *Session) error {
	session.sessionDir = me.constructSessionDir(session.id)
	return me.ensureDir(session.sessionDir)
}

func (me *Manager) constructSessionDir(key string) string {
	return filepath.Join(me.sessionsDir, "dir"+key)
}

func (me *Manager) Clean() error {
	return me.sessionsDB.Clean()
}

func (me *Manager) Close() (err error) {
	me.sessionsDB.IterateMemStorage(func(key string, val interface{}) {
		if sess, ok := val.(*Session); ok {
			err = sess.close()
			if err != nil {
				log.Println(err.Error())
			}
		}
	})
	me.sessionsDB.Close()
	return nil
}
