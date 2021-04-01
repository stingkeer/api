package call

import (
	"gitee.com/fast_api/api/public"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime"
)

type callerDefault struct {
	serialize public.Serialize
}

var adapters = make(map[reflect.Type]Adapter)

func NewCaller(serialize public.Serialize) *callerDefault {
	return &callerDefault{
		serialize: serialize,
	}
}

func RegisterTypeMapper(adapter Adapter) {
	if adapter != nil {
		for _, m := range adapter.Register() {
			adapters[m] = adapter
		}
	}
}

//request == call(def) => value
func (c *callerDefault) Call(f *public.Entry, req *http.Request) interface{} {
	v := reflect.ValueOf(f.Fn)
	if v.Type().NumIn() > 1 {
		logrus.Error("not support param > 1")
		return nil
	}
	name := runtime.FuncForPC(reflect.ValueOf(f.Fn).Pointer()).Name()
	m := c.getFuncInfo(name)
	if m == nil {
		logrus.Error("not find method in header")
		os.Exit(2)
	}
	params := req.URL.Query()
	for k, v := range f.Ids {
		params.Add(k, v)
	}
	var paramsV []reflect.Value
	for name, p := range m.Param {
		pw := public.ParamWarp{Request: *req}
		pw.PTyp = v.Type().In(p.Order)
		pw.PName = name
		if t, b := adapters[p.Typ]; b {
			if v, b := params[name]; b {
				pw.PValue = v[0]
			}
			paramsV = append(paramsV, t.Mapper(pw))
		} else if pw.PTyp.Kind() == reflect.Struct && req.Method == http.MethodPost {
			newT := reflect.New(pw.PTyp)
			bytes, _ := ioutil.ReadAll(req.Body)
			err := c.serialize.Decode(bytes, newT.Interface())
			if err != nil {
				panic(err)
			}
			paramsV = append(paramsV, newT.Elem())
		} else { //default value
			logrus.Tracef("not support %s set default value", pw.PTyp.Kind())
			paramsV = append(paramsV, c.defaultCallValue(pw.PTyp.Kind()))
		}
	}
	vs := v.Call(paramsV)
	if len(vs) == 0 {
		logrus.Warn("call method no return")
		return nil
	}
	return vs[0].Interface()
}

func toPtr(obj interface{}) reflect.Value {
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	return vp
}

func (c *callerDefault) getFuncInfo(name string) *public.MethodInfo {
	if m, ok := public.MethodsPools[name]; ok {
		return &m
	}
	logrus.Errorf("not find name [%s]", name)
	return nil
}

//other param set default value
func (c *callerDefault) defaultCallValue(kind reflect.Kind) reflect.Value {
	switch kind {
	case reflect.String:
		return reflect.ValueOf("")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(0)
	}
	return reflect.ValueOf(nil)
}

func init() {
	RegisterTypeMapper(&BaseType{})
	RegisterTypeMapper(&bigType{})
	RegisterTypeMapper(&FileType{})
	RegisterTypeMapper(&HttpType{})
	RegisterTypeMapper(&headType{})
}
