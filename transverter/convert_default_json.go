package transverter

import (
	"encoding/json"
	"gitee.com/fast_api/api/public"
	"reflect"
)

type JSONConvertImpl struct {
}

func (c *JSONConvertImpl) ConvertFrom(bytes []byte, vpr interface{}) error {
	return json.Unmarshal(bytes, vpr)
}

func (c *JSONConvertImpl) ConvertTo(f interface{}) *public.Header {
	var head public.Header
	kind := reflect.Indirect(reflect.ValueOf(f)).Kind()
	if kind == reflect.String {
		head.Bytes = []byte(f.(string))
		head.ContentType = "text/plain"
		return &head
	}
	bytes, _ := json.Marshal(f)
	head.Bytes = bytes
	head.ContentType = "application/json"
	return &head
}
