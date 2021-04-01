package call

import (
	"gitee.com/fast_api/api/public"
	"reflect"
)

type Adapter interface {
	Mapper(param public.ParamWarp) reflect.Value
	Register() []reflect.Type
}
