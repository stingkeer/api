package intercept

import (
	"net/http"
	"reflect"

	"gitee.com/fast_api/api/def"
)

type HttpIntercept interface {
	// HttpIntercept
	// if return true, indicating interception and not executing backwards
	Http(rw http.ResponseWriter, req *http.Request) bool
	Order() def.HandlerOrder
}

type MethodIntercept interface {
	Invoke(m *def.MethodInfo, args []reflect.Value) []reflect.Value
}
