package types

import (
	"fmt"
	"gitee.com/fast_api/api/def"
	"reflect"
)

type TypeRequire struct {
	BaseType
}

func (b *TypeRequire) Mapper(p def.ParamWarp) reflect.Value {
	if p.PValue == "" {
		panic(fmt.Sprintf("param %s is require", p.PName))
	}
	return b.BaseType.Mapper(p).Convert(p.PTyp)
}

func (b *TypeRequire) Register() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*def.Int8Req)(nil)).Elem(),
		reflect.TypeOf((*def.Int16Req)(nil)).Elem(),
		reflect.TypeOf((*def.Int32Req)(nil)).Elem(),
		reflect.TypeOf((*def.Int64Req)(nil)).Elem(),
		reflect.TypeOf((*def.StringReq)(nil)).Elem(),
	}
}
