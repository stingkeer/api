package cache

import (
	"bytes"
	"fmt"
	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/mg"
	"io"
	"reflect"
	"strings"
	"time"
)

var (
	cacheType                         = reflect.TypeOf((*Cache)(nil)).Elem()
	persistenceCache PersistenceCache = &defaultPersistenceCache{cache: make(map[string]cEntry)}
	processCache     ProcessCache     = &defaultProcessCacheImpl{}
)

func init() {
	http.RegisterReturnHandler(&Bytes{})
	call.SetMethodProxy(func(fn call.MethodCaller, m *def.MethodInfo, args []reflect.Value) []reflect.Value {
		var builder strings.Builder
		for _, dm := range m.Param {
			builder.WriteString(dm.Name)
			builder.WriteString("=")
			builder.WriteString(fmt.Sprintf("%s", args[dm.Order].Interface()))
			builder.WriteString("@")
		}
		key := processCache.EncodeKey(m.MethodName, builder.String())
		if cv := persistenceCache.Get(key); cv != nil {
			return []reflect.Value{
				reflect.ValueOf(Bytes(cv)),
			}
		}
		invokeReturn := fn.Invoke(m, args)
		if len(invokeReturn) > 1 {
			cache := validCache(invokeReturn)
			if cache != nil {
				mg.Invoke(func(serialize def.Serialize) {
					persistenceCache.Set(key, processCache.EncodeValue(serialize, invokeReturn[0]), cache.ExpireTime_())
				})
			}
		}
		return invokeReturn
	})

}

func validArgs(m *def.MethodInfo, args []reflect.Value) {

}

func validCache(vs []reflect.Value) Cache {
	for _, v := range vs {
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
