package dwarf

import (
	"fmt"
	"reflect"
)

type ArgsMeta struct {
	Order int
	Name  string
	Typ   reflect.Type
}

func (m *ArgsMeta) String() string {
	return fmt.Sprintf("%s %s", m.Name, m.Typ)
}

type MethodMeta struct {
	MethodName string
	Args       []ArgsMeta
	Ret        []ArgsMeta
}

func (m *MethodMeta) String() string {
	return fmt.Sprintf("%s(%#v)", m.MethodName, m.Args)
}
