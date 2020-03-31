package api

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
)

type CallerDefault struct {
	convert Convert
}

func (c *CallerDefault) call(f interface{}, req *http.Request) interface{} {
	switch req.Method {
	case "GET":
		return c.callGet(f, req)
	case "POST":
		return c.callPost(f, req)
	}
	return nil
}

func (c *CallerDefault) callPost(f interface{}, req *http.Request) interface{} {
	v := reflect.ValueOf(f)
	newT := reflect.New(v.Type().In(0))
	bytes, _ := ioutil.ReadAll(req.Body)
	c.convert.convertFrom(bytes, newT.Interface())
	vs := v.Call([]reflect.Value{newT.Elem()})
	return vs[0].Interface()
}

func (c *CallerDefault) callGet(f interface{}, req *http.Request) interface{} {
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
	params := req.URL.Query()
	for name, p := range m.Param {
		if v, b := params[name]; b {
			pvs[p.Order] = c.typeConvert(v[0], t.In(p.Order))
		} else {
			pvs[p.Order] = c.defaultCallValue(t.In(p.Order).Kind())
		}
	}
	vs := reflect.ValueOf(f).Call(pvs)
	return vs[0].Interface()
}

func (c *CallerDefault) typeConvert(value string, dest reflect.Type) reflect.Value {
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
	default:

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
