package types

import (
	"fmt"
	"math/big"
	"reflect"

	"gitee.com/fast_api/api/def"
)

var _ def.Adapter = (*BigType)(nil)

type BigType struct {
}

func (f BigType) Mapper(param *def.ParamWarp) reflect.Value {
	if param.PTyp.Kind() == reflect.Struct {
		if in, b := new(big.Int).SetString(param.PValue, 10); b {
			return reflect.ValueOf(*in).Convert(param.PTyp)
		}
	}
	if param.PTyp.Kind() == reflect.Ptr {
		if in, b := new(big.Int).SetString(param.PValue, 10); b {
			return reflect.ValueOf(in).Convert(param.PTyp)
		}
	}
	panic(fmt.Sprintf("BigType error param %s", param.PValue))

}

func (f BigType) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*big.Int)(nil)).Elem(), reflect.TypeOf((*big.Int)(nil))}
}
