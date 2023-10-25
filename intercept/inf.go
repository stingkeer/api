package intercept

import (
	"net/http"
	"reflect"

	"gitee.com/fast_api/api/def"
)

type HttpIntercept interface {
	Http(rw http.ResponseWriter, req *http.Request) bool
	Order() def.HandlerOrder
}

type MethodIntercept interface {
	Invoke(m *def.MethodInfo, args []reflect.Value) []reflect.Value
}
