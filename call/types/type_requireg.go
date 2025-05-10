package types

import (
	"reflect"

	"go.aew.app/api.v1/def"
)

var _ def.Adapter = (*TypeRequireG)(nil)

type TypeRequireG struct {
	TypeRequire
}

func (b *TypeRequireG) Mapper(p *def.ParamWarp) reflect.Value {
	return b.TypeRequire.Mapper(p)
}

func (b *TypeRequireG) Register() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*def.Int[any])(nil)).Elem(),
		reflect.TypeOf((*def.Int8[any])(nil)).Elem(),
		reflect.TypeOf((*def.Int16[any])(nil)).Elem(),
		reflect.TypeOf((*def.Int32[any])(nil)).Elem(),
		reflect.TypeOf((*def.Int64[any])(nil)).Elem(),
		reflect.TypeOf((*def.String[any])(nil)).Elem(),
	}
}
