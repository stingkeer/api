package serialize

import (
	"encoding/json"
	"gitee.com/fast_api/api/def"
	"reflect"
)

type JsonConvertImpl struct{}

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
		bytes, e := json.Marshal(f)
		if e != nil {
			panic(e)
		}
		ctxt.Bytes = bytes
		return &ctxt
	}
}
