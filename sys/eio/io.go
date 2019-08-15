package eio

import (
	"reflect"
	"sync"
)

type (
	Params struct {
		Args   []Arg
		Return []Arg
	}
	Arg struct {
		Kind reflect.Kind
		Name string
	}
	IO interface {
		Init(obj map[string]interface{}) (IO, error)
		TestInit(obj map[string]interface{}) (IO, error)
		Call(methodName string, args ...interface{}) (interface{}, error)
		TestCall(methodName string, args ...interface{}) (interface{}, error)
		ListInitArgs() []Arg
		ListMethodNames() []string
		ListParams(methodName string) Params
		Close() error
	}

	IOManager struct {
		providers      map[string]IO
		providersCache []string
		providersLock  sync.RWMutex
	}
)

func NewIOManager() (*IOManager, error) {
	return &IOManager{providers: make(map[string]IO)}, nil
}

func (me *IOManager) PlugIn(providerName string, pio IO) {
	me.providersLock.Lock()
	me.providers[providerName] = pio
	me.updateCache()
	me.providersLock.Unlock()
}

func (me *IOManager) Get(providerName string) IO {
	me.providersLock.Lock()
	defer me.providersLock.Unlock()
	return me.providers[providerName]
}

func (me *IOManager) PlugOut(providerName string) {
	me.providersLock.Lock()
	delete(me.providers, providerName)
	me.updateCache()
	me.providersLock.Unlock()
}

func (me *IOManager) updateCache() {
	me.providersCache = make([]string, len(me.providers))
	i := 0
	for k := range me.providers {
		me.providersCache[i] = k
		i++
	}
}

func (me *IOManager) ListProviders() []string {
	me.providersLock.RLock()
	defer me.providersLock.RUnlock()
	return me.providersCache
}
