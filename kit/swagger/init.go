package swagger

import (
	"mime/multipart"
	"reflect"

	"go.aew.app/api.v1/call/types"
)

// passTypes are framework-injected parameter types (http.Request, def.Header,
// *ws.WSCtx). They never appear in request input and are skipped by the doc
// generator. big.Int and multipart.Reader are deliberately NOT here: they are
// real wire parameters (integer query / file upload) and must be documented.
var passTypes []reflect.Type

var multipartReaderType = reflect.TypeOf((*multipart.Reader)(nil)).Elem()

func init() {
	passTypes = append(passTypes, types.HttpType{}.Register()...)
	passTypes = append(passTypes, types.HeadType{}.Register()...)
	passTypes = append(passTypes, (&types.WSType{}).Register()...)
}

// isPassType reports whether t (or its pointer element) is an injected
// infrastructure type.
func isPassType(t reflect.Type) bool {
	for _, e := range passTypes {
		if t == e {
			return true
		}
		if t.Kind() == reflect.Ptr && t.Elem() == e {
			return true
		}
		if e.Kind() == reflect.Ptr && t == e.Elem() {
			return true
		}
	}
	return false
}
