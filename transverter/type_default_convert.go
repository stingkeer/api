package transverter

import (
	"github.com/sirupsen/logrus"
	"math/big"
	"reflect"
)

type DefaultTypeConvert struct {
}

func (cvr *DefaultTypeConvert) ConvertTo(value string, typ reflect.Type) reflect.Value {
	switch {
	case typ == reflect.TypeOf((*big.Int)(nil)).Elem():
		in, _ := new(big.Int).SetString(value, 10)
		return reflect.ValueOf(*in).Convert(typ)
	default:
		logrus.Error("not support")
	}
	return reflect.Value{}
}
