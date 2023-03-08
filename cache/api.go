package cache

import (
	"gitee.com/fast_api/api/def"
	"reflect"
	"strings"
	"time"
)

type ProcessCache interface {
	EncodeKey(v ...any) []byte
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

func (d *defaultProcessCacheImpl) EncodeKey(v ...any) []byte {
	var s strings.Builder
	for _, a := range v {
		s.WriteString(a.(string))
	}
	return []byte(s.String())
}

func (d *defaultProcessCacheImpl) EncodeValue(s def.Serialize, v reflect.Value) []byte {
	return s.Encode(v.Interface()).Bytes
}

type cEntry struct {
	data []byte
	t    time.Time
}
type defaultPersistenceCache struct {
	cache map[string]cEntry
}

func (d defaultPersistenceCache) Set(key []byte, value []byte, ttl time.Duration) {
	d.cache[string(key)] = cEntry{
		data: value,
		t:    time.Now().Add(ttl),
	}
}

func (d defaultPersistenceCache) Get(key []byte) []byte {
	if f, b := d.cache[string(key)]; b && f.t.After(time.Now()) {
		return f.data
	}
	return nil
}
