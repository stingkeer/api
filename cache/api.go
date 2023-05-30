package cache

import (
	"fmt"
	"gitee.com/fast_api/api/def"
	"reflect"
	"strings"
	"sync"
	"time"
)

type ProcessCache interface {
	EncodeKey(m *def.MethodInfo, args []reflect.Value) []byte
	EncodeValue(s def.Serialize, v reflect.Value) []byte
}

type PersistenceCache interface {
	Set(key []byte, value []byte, ttl time.Duration)
	Get(key []byte) []byte
}

func SetPersistenceCacheImpl(persistenceCacheI PersistenceCache) {
	persistenceCache = persistenceCacheI
}

func SetProcessCacheImpl(processCacheI ProcessCache) {
	processCache = processCacheI
}

type defaultProcessCacheImpl struct{}

func (d *defaultProcessCacheImpl) EncodeKey(m *def.MethodInfo, args []reflect.Value) []byte {
	var builder strings.Builder
	for _, dm := range m.Param {
		if !args[dm.Order].IsValid() {
			continue
		}
		builder.WriteString(dm.Name)
		builder.WriteString("=")
		builder.WriteString(fmt.Sprintf("%s", args[dm.Order].Interface()))
		builder.WriteString("@")
	}
	return []byte(builder.String())
}

func (d *defaultProcessCacheImpl) EncodeValue(s def.Serialize, v reflect.Value) []byte {
	return s.Encode(v.Interface()).Bytes
}

type cEntry struct {
	data []byte
	t    time.Time
}
type defaultPersistenceCache struct {
	//cache map[string]cEntry
	cache sync.Map
}

func (d defaultPersistenceCache) Set(key []byte, value []byte, ttl time.Duration) {
	d.cache.Store(string(key), cEntry{
		data: value,
		t:    time.Now().Add(ttl),
	})
}

func (d defaultPersistenceCache) Get(key []byte) []byte {
	if f, b := d.cache.Load(string(key)); b && f != nil {
		entry := f.(cEntry)
		if entry.t.After(time.Now()) {
			return entry.data
		}
	}
	return nil
}
