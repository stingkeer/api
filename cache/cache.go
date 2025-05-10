package cache

import (
	"bytes"
	"io"
	"reflect"
	"time"

	"go.aew.app/api.v1/call"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/http"
)

var (
	cacheType                         = reflect.TypeOf((*Cache)(nil)).Elem()
	persistenceCache PersistenceCache = &defaultPersistenceCache{}
	processCache     ProcessCache     = &defaultProcessCacheImpl{}
)

func init() {
	http.RegisterReturnHandler(&Bytes{})
	call.SetMethodProxy(func(fn call.MethodCaller, m *def.MethodInfo, args []reflect.Value) []reflect.Value {
		key := []byte(m.MethodName + "@")
		encodeKey := processCache.EncodeKey(m, args)
		if len(encodeKey) > 0 {
			key = append(key, encodeKey...)
		}
		if cv := persistenceCache.Get(key); cv != nil {
			return []reflect.Value{
				reflect.ValueOf(Bytes(cv)),
			}
		}
		invokeReturn := fn.Invoke(m, args)
		if len(invokeReturn) > 1 {
			cache := validCache(invokeReturn)
			if cache != nil {
				persistenceCache.Set(key, processCache.EncodeValue(def.DefaultContext.Serialize, invokeReturn[0]), cache.ExpireTime_())
			}
		}
		return invokeReturn
	})

}

func validArgs(m *def.MethodInfo, args []reflect.Value) {

}

func validCache(vs []reflect.Value) Cache {
	for _, v := range vs {
		if !v.IsValid() || v.IsNil() {
			continue
		}
		if v.Type().Implements(cacheType) {
			return v.Interface().(Cache)
		}
	}
	return nil
}

type Bytes []byte

func (b Bytes) Return() io.Reader {
	return bytes.NewReader(b)
}
func (b Bytes) Register() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Bytes)(nil)).Elem(),
	}
}

func (b Bytes) ContentType() string {
	return def.Content_JSON
}

type Cache interface {
	ExpireTime_() time.Duration
}

type Key interface {
}

type ImplCache struct {
	time time.Duration
}

func NewCacheImpl(time time.Duration) *ImplCache {
	return &ImplCache{time: time}
}

func (c *ImplCache) ExpireTime_() time.Duration {
	return c.time
}
