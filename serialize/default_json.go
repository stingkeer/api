package serialize

import (
	"encoding/json"
	"gitee.com/fast_api/api/public"
	"reflect"
)

type JsonConvertImpl struct{}

func (c *JsonConvertImpl) Decode(bytes []byte, vpr interface{}) error {
	return json.Unmarshal(bytes, vpr)
}

func (c *JsonConvertImpl) Encode(f interface{}) *public.Content {
	var ctxt public.Content
	ctxt.ContentType = public.Content_JSON
	if e, b := f.(error); b {
		panic(e)
	}
	kind := reflect.Indirect(reflect.ValueOf(f)).Kind()
	switch kind {
	case reflect.String:
		ctxt.Bytes = []byte(f.(string))
		return &ctxt
	default:
		bytes, _ := json.Marshal(f)
		ctxt.Bytes = bytes
		return &ctxt
	}
}
