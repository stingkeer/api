package api

import (
	"encoding/json"
	"reflect"
)

type JSONConvertImpl struct {
	contentType string
}

func (c *JSONConvertImpl) convert(f interface{}) []byte {
	kind := reflect.Indirect(reflect.ValueOf(f)).Kind()
	if kind == reflect.String {
		c.contentType = "text/plain"
		return []byte(f.(string))
	}
	bytes, _ := json.Marshal(f)
	c.contentType = "application/json"
	return bytes
}

func (c *JSONConvertImpl) getContentType() string {
	return c.contentType
}
