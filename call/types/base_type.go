package types

import (
	"reflect"
	"strconv"

	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/log"
	"go.aew.app/api.v1/utils"
)

var _ def.Adapter = (*BaseType)(nil)

type BaseType struct{}

func (b *BaseType) Mapper(p *def.ParamWarp) reflect.Value {
	dest := p.PTyp
	value := p.PValue
	if value == "" {
		return utils.DefaultCallValue(dest)
	}
	switch dest.Kind() {
	case reflect.Bool:
		parseBool, err := strconv.ParseBool(value)
		if err != nil {
			panic(err)
		}
		return reflect.ValueOf(parseBool)
	case reflect.String:
		return reflect.ValueOf(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s, e := strconv.ParseInt(value, 10, 64)
		if e != nil {
			panic(e)
		}
		return reflect.ValueOf(s).Convert(dest)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s, e := strconv.ParseUint(value, 10, 64)
		if e != nil {
			panic(e)
		}
		return reflect.ValueOf(s).Convert(dest)
	case reflect.Float32, reflect.Float64:
		s, e := strconv.ParseFloat(value, 64)
		if e != nil {
			panic(e)
		}
		return reflect.ValueOf(s).Convert(dest)
	default:
		log.Errorf("not find type %s", dest)

	}
	return reflect.ValueOf(nil)
}

func (b *BaseType) Register() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*bool)(nil)).Elem(),
		reflect.TypeOf((*string)(nil)).Elem(),
		reflect.TypeOf((*int)(nil)).Elem(),
		reflect.TypeOf((*int8)(nil)).Elem(),
		reflect.TypeOf((*int16)(nil)).Elem(),
		reflect.TypeOf((*int32)(nil)).Elem(),
		reflect.TypeOf((*int64)(nil)).Elem(),
		reflect.TypeOf((*uint)(nil)).Elem(),
		reflect.TypeOf((*uint8)(nil)).Elem(),
		reflect.TypeOf((*uint16)(nil)).Elem(),
		reflect.TypeOf((*uint32)(nil)).Elem(),
		reflect.TypeOf((*uint64)(nil)).Elem(),
		reflect.TypeOf((*float32)(nil)).Elem(),
		reflect.TypeOf((*float64)(nil)).Elem(),
	}
}
