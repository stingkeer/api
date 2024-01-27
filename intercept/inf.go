package intercept

import (
	"net/http"
	"reflect"
	"sync"

	"gitee.com/fast_api/api/def"
)

type HttpIntercept interface {
	// HttpIntercept
	// if return true, indicating interception and not executing backwards
	Http(rw http.ResponseWriter, req *http.Request, ctx *HttpContext) bool
	Order() def.HandlerOrder
}

type HttpContext struct {
	sync.Map
}

func NewHttpContext() *HttpContext {
	return &HttpContext{}
}

func (hc *HttpContext) SkipResponse() {
	hc.Store("SkipResponse", 1)
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
