package call

import (
	"gitee.com/fast_api/api/public"
	"math/big"
	"reflect"
)

type bigType struct {
}

func (f bigType) Mapper(param public.ParamWarp) reflect.Value {
	in, _ := new(big.Int).SetString(param.PValue, 10)
	return reflect.ValueOf(*in).Convert(param.PTyp)
}

func (f bigType) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*big.Int)(nil)).Elem()}
}
