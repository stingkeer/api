package types

import (
	"fmt"
	"reflect"

	"go.aew.app/api.v1/def"
)

var _ def.Adapter = (*TypeRequire)(nil)

type TypeRequire struct {
	BaseType
}

func (b *TypeRequire) Mapper(p *def.ParamWarp) reflect.Value {
	if p.PValue == "" {
		panic(fmt.Sprintf("param %s is require", p.PName))
	}
	pv := reflect.New(p.PTyp)
	field0 := pv.Elem().Field(0)
	p.PTyp = field0.Type()
	field0.Set(b.BaseType.Mapper(p).Convert(p.PTyp))
	return pv.Elem()
}

func (b *TypeRequire) Register() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*def.IntReq)(nil)).Elem(),
		reflect.TypeOf((*def.Int8Req)(nil)).Elem(),
		reflect.TypeOf((*def.Int16Req)(nil)).Elem(),
		reflect.TypeOf((*def.Int32Req)(nil)).Elem(),
		reflect.TypeOf((*def.Int64Req)(nil)).Elem(),
		reflect.TypeOf((*def.StringReq)(nil)).Elem(),
	}
}
