package swagger

import (
	"reflect"

	"go.aew.app/api/call/types"
)

var (
	defTypes []reflect.Type
)

func init() {
	defTypes = append(defTypes, types.HttpType{}.Register()...)
	defTypes = append(defTypes, types.HeadType{}.Register()...)
	defTypes = append(defTypes, (&types.WSType{}).Register()...)
	defTypes = append(defTypes, (types.BigType{}).Register()...)
	defTypes = append(defTypes, (types.FileType{}).Register()...)
}

func isDefTypes(t reflect.Type) bool {
	for _, e := range defTypes {
		if e == t {
			return true
		}
	}
	return false
}
