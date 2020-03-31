package api

import (
	"encoding/json"
	"reflect"
)

type JSONConvertImpl struct {
}

func (c *JSONConvertImpl) convertFrom(bytes []byte, vpr interface{}) error {
	return json.Unmarshal(bytes, vpr)
}

func (c *JSONConvertImpl) convertTo(f interface{}) *Header {
	var head Header
	kind := reflect.Indirect(reflect.ValueOf(f)).Kind()
	if kind == reflect.String {
		head.bytes = []byte(f.(string))
		head.ContentType = "text/plain"
		return &head
	}
	bytes, _ := json.Marshal(f)
	head.bytes = bytes
	head.ContentType = "application/json"
	return &head
}
