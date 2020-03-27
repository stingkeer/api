package api

import (
	"net/url"
	"reflect"
)

type CallerDefault struct {
}

func (c *CallerDefault) call(f interface{}, params url.Values) interface{} {
	vs := reflect.ValueOf(f).Call(nil)
	return vs[0].Interface()
}
