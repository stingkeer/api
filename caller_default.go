package api

import (
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strings"
)

type CallerDefault struct {
}

func (c *CallerDefault) call(f interface{}, params url.Values) interface{} {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	m := c.getFuncInfo(name)
	t := reflect.TypeOf(f)
	var pvs = make([]reflect.Value, t.NumIn())
	logrus.Tracef("method has param [%d]", t.NumIn())
	for name, p := range m.Param {
		if v, b := params[name]; b {
			pvs[p.Order] = reflect.ValueOf(v[0])
		} else {
			pvs[p.Order] = c.defaultCallValue(t.In(p.Order).Kind())
		}
	}
	vs := reflect.ValueOf(f).Call(pvs)
	return vs[0].Interface()
}

func (c *CallerDefault) getFuncInfo(name string) *MethodInfo {
	for _, method := range methods {
		if strings.HasSuffix(name, method.MethodName) {
			return &method
		}
	}
	logrus.Errorf("not find name %s", name)
	os.Exit(2)
	return nil
}

func (c *CallerDefault) defaultCallValue(kind reflect.Kind) reflect.Value {
	switch kind {
	case reflect.String:
		return reflect.ValueOf("")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(0)
	}
	return reflect.ValueOf(nil)
}
