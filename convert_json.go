package api

import (
	"encoding/json"
	"reflect"
)

type JSONConvertImpl struct {
	contentType string
}

func (c *JSONConvertImpl) convert(f interface{}) *Header {
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

func (c *JSONConvertImpl) getContentType() string {
	return c.contentType
}
