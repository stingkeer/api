package types

import (
	"fmt"
	"gitee.com/fast_api/api/def"
	"math/big"
	"reflect"
)

type BigType struct {
}

func (f BigType) Mapper(param def.ParamWarp) reflect.Value {
	if in, b := new(big.Int).SetString(param.PValue, 10); b {
		return reflect.ValueOf(*in).Convert(param.PTyp)
	}
	panic(fmt.Sprintf("BigType error param %s", param.PValue))

}

func (f BigType) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*big.Int)(nil)).Elem()}
}
