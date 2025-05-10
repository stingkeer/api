package http

import (
	"reflect"

	"go.aew.app/api/def"
)

var retAdapters = make(map[reflect.Type]def.RetAdapter)

func RegisterReturnHandler(ret def.RetAdapter) {
	if retAdapters != nil {
		for _, m := range ret.Register() {
			retAdapters[m] = ret
		}
	}
}
