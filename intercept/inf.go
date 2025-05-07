package intercept

import (
	"net/http"
	"reflect"

	"gitee.com/fast_api/api/def"
)

type HttpIntercept interface {
	// HttpIntercept
	// if return true, indicating interception and not executing backwards
	Http(rw http.ResponseWriter, req *http.Request, ctx *HttpContext) bool
	Order() def.HandlerOrder
}

type HttpContext struct {
	m map[string]any
}

func NewHttpContext() *HttpContext {
	return &HttpContext{m: make(map[string]any)}
}

func (hc *HttpContext) Clear() {
	clear(hc.m)
}

func (hc *HttpContext) Store(key string, v any) {
	hc.m[key] = v
}

func (hc *HttpContext) LoadAndDelete(key string) (any, bool) {
	if v, b := hc.m[key]; b {
		delete(hc.m, key)
		return v, b
	}
	return nil, false
}

// SkipResponse
// Skip Response
func (hc *HttpContext) SkipResponse() {
	hc.Store("SkipResponse", 1)
}

func (hc *HttpContext) Load(key string) (any, bool) {
	if v, b := hc.m[key]; b {
		return v, b
	} else {
		return nil, false
	}

}

func (hc *HttpContext) IsSkipResponse() bool {
	if v, b := hc.Load("SkipResponse"); b && v == 1 {
		return true
	}
	return false
}

type MethodIntercept interface {
	Invoke(m *def.MethodInfo, args []reflect.Value) []reflect.Value
}
