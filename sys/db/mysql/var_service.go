package mysql

import (
	"database/sql"
	"sync"
)

type (
	varNotificationConsumer interface {
		syncVarsTask(something interface{})
	}
	varService struct {
		db       *sql.DB
		mutex    sync.Mutex
		running  bool
		consumer varNotificationConsumer
		notify   chan interface{}
		closed   bool
	}
)

func newVarService(db *sql.DB, consumer varNotificationConsumer) *varService {
	return &varService{db: db, consumer: consumer, notify: make(chan interface{}, 100)}
}

func (me *varService) notifyUpdate(something interface{}) {
	me.mutex.Lock()
	defer me.mutex.Unlock()
	if !me.closed {
		me.notify <- something
		if !me.running {
			me.running = true
			go me.syncedThread()
		}
	}
}

func (me *varService) syncedThread() {
	defer func() {
		me.running = false
		me.consumer = nil
		me.db = nil
	}()
	for {
		select {
		case obj, ok := <-me.notify:
			if !ok {
				return
			} //channel closed
			if obj != nil {
				me.consumer.syncVarsTask(obj)
			}
		}
	}
}

func (me *varService) Close() {
	me.mutex.Lock()
	defer me.mutex.Unlock()
	if me.notify != nil {
		close(me.notify)
	}
	me.closed = true
}
