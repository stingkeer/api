package http

import (
	"gitee.com/fast_api/api/def"
	"reflect"
)

var retAdapters = make(map[reflect.Type]def.RetAdapter)

func RegisterReturnHandler(ret def.RetAdapter) {
	if retAdapters != nil {
		for _, m := range ret.Register() {
			retAdapters[m] = ret
		}
	}
}
