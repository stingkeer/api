package api

import "reflect"

type CallerDefault struct {
}

func (c *CallerDefault) call(f interface{}) interface{} {
	vs := reflect.ValueOf(f).Call(nil)
	return vs[0].Interface()
}
