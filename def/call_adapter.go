package def

import (
	"io"
	"reflect"
)

type Adapter interface {
	Mapper(param ParamWarp) reflect.Value
	Register() []reflect.Type
}

type RetAdapter interface {
	ContentType
	Return() io.Reader
	Register() []reflect.Type
}
