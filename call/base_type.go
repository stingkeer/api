package call

import (
	"gitee.com/fast_api/api/public"
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
)

type BaseType struct{}

func getFuncInfo(name string) *public.MethodInfo {
	if m, ok := public.MethodsPools[name]; ok {
		return &m
	}
	logrus.Errorf("not find name [%s]", name)
	return nil
}

func (b *BaseType) Mapper(p public.ParamWarp) reflect.Value {
	dest := p.PTyp
	value := p.PValue
	switch dest.Kind() {
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
		s, e := strconv.ParseFloat(value, 10)
		if e != nil {
			panic(e)
		}
		return reflect.ValueOf(s).Convert(dest)
	default:
		logrus.Errorf("not find type %s", dest)

	}
	return reflect.ValueOf(nil)
}

func (b *BaseType) Register() []reflect.Type {
	return []reflect.Type{
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
