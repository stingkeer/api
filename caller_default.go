package api

import (
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strconv"
)

type CallerDefault struct {
}

func (c *CallerDefault) call(f interface{}, params url.Values) interface{} {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	logrus.Tracef("call function name [%s]", name)
	m := c.getFuncInfo(name)
	if m == nil {
		logrus.Error("not find method in header")
		os.Exit(2)
	}
	t := reflect.TypeOf(f)
	var pvs = make([]reflect.Value, t.NumIn())
	logrus.Tracef("method has param [%d]", t.NumIn())
	for name, p := range m.Param {
		if v, b := params[name]; b {
			pvs[p.Order] = c.convert(v[0], t.In(p.Order))
		} else {
			pvs[p.Order] = c.defaultCallValue(t.In(p.Order).Kind())
		}
	}
	vs := reflect.ValueOf(f).Call(pvs)
	return vs[0].Interface()
}

func (c *CallerDefault) convert(value string, dest reflect.Type) reflect.Value {
	switch dest.Kind() {
	case reflect.String:
		return reflect.ValueOf(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s, e := strconv.ParseInt(value, 10, 64)
		if e != nil {
			panic(e)
		}
		return reflect.ValueOf(s).Convert(dest)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s, e := strconv.ParseUint(value, 10, 64)
		if e != nil {
			panic(e)
		}
		return reflect.ValueOf(s).Convert(dest)
	case reflect.Float32, reflect.Float64:
		s, e := strconv.ParseFloat(value, 10)
		if e != nil {
			panic(e)
		}
		return reflect.ValueOf(s).Convert(dest)
	}
	return reflect.ValueOf(nil)
}

func (c *CallerDefault) getFuncInfo(name string) *MethodInfo {
	sStruct, sFunc := SplitFuncName(name)
	for _, method := range methods {
		if sStruct == "" && sFunc == method.MethodName {
			return &method
		}
		if method.Receive == sStruct && sFunc == method.MethodName {
			return &method
		}
	}
	logrus.Errorf("not find name [%s]", name)
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
