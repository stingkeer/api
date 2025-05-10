package serialize

import (
	"encoding/json"
	"reflect"

	"go.aew.app/api.v1/def"
)

var _ def.Serialize = (*JsonConvertImpl)(nil)

type JsonConvertImpl struct{}

// ContentType implements def.Serialize.
func (c *JsonConvertImpl) ContentType() string {
	return def.Content_JSON
}

func (c *JsonConvertImpl) Decode(bytes []byte, vpr interface{}) error {
	return json.Unmarshal(bytes, vpr)
}

func (c *JsonConvertImpl) Encode(f interface{}) *def.Content {
	var ctxt def.Content
	ctxt.ContentType = def.Content_JSON
	if e, b := f.(error); b {
		panic(e)
	}
	kind := reflect.Indirect(reflect.ValueOf(f)).Kind()
	switch kind {
	case reflect.String:
		ctxt.Bytes = []byte(f.(string))
		return &ctxt
	default:
		bytes, err := json.Marshal(f)
		if err != nil {
			panic(err)
		}
		ctxt.Bytes = bytes
		return &ctxt
	}
}
