package intercept

import (
	"gitee.com/fast_api/api/def"
	"net/http"
	"reflect"
)

type HttpIntercept interface {
	Http(rw http.ResponseWriter, req *http.Request) bool
	Order() int
}

type MethodIntercept interface {
	Invoke(fn reflect.Value, m *def.MethodInfo, args []reflect.Value) []reflect.Value
}
