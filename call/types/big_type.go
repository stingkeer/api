package types

import (
	"gitee.com/fast_api/api/def"
	"math/big"
	"reflect"
)

type BigType struct {
}

func (f BigType) Mapper(param def.ParamWarp) reflect.Value {
	in, _ := new(big.Int).SetString(param.PValue, 10)
	return reflect.ValueOf(*in).Convert(param.PTyp)
}

func (f BigType) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*big.Int)(nil)).Elem()}
}
