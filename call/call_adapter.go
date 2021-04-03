package call

import "reflect"

type RealCall interface {
	Invoke(value reflect.Value) interface{}
}
