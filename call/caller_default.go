package call

import (
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/log"
	"gitee.com/fast_api/api/utils"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime"
)

type callerDefault struct {
	serialize def.Serialize
}

var (
	adapters = make(map[reflect.Type]def.Adapter)
)

func NewCaller(serialize def.Serialize) *callerDefault {
	return &callerDefault{
		serialize: serialize,
	}
}

func RegisterTypeMapper(adapter def.Adapter) {
	if adapter != nil {
		for _, m := range adapter.Register() {
			adapters[m] = adapter
		}
	}
}

// Call request == call(def) => value
func (c *callerDefault) Call(f *def.Entry, req *http.Request) interface{} {
	v := reflect.ValueOf(f.Fn)
	name := runtime.FuncForPC(reflect.ValueOf(f.Fn).Pointer()).Name()
	m := c.getFuncInfo(name)
	if m == nil {
		log.Error("not find method in exe")
		os.Exit(2)
	}
	params := req.URL.Query()
	for k, v := range f.Ids {
		params.Add(k, v)
	}

	paramsV := make([]reflect.Value, len(m.Param))

	for name, p := range m.Param {
		pw := def.ParamWarp{Request: *req}
		pw.PTyp = v.Type().In(p.Order)
		pw.PName = name
		if t, b := adapters[p.Typ]; b {
			if v, b := params[name]; b {
				pw.PValue = v[0]
			}
			paramsV[p.Order] = t.Mapper(pw)
		} else if pw.PTyp.Kind() == reflect.Struct && req.Method == http.MethodPost {
			newT := reflect.New(pw.PTyp)
			bytes, _ := ioutil.ReadAll(req.Body)
			err := c.serialize.Decode(bytes, newT.Interface())
			if err != nil {
				panic(err)
			}
			paramsV[p.Order] = newT.Elem()
		} else { //default value
			log.Tracef("not support %s set default value", pw.PTyp)
			paramsV[p.Order] = utils.DefaultCallValue(pw.PTyp.Kind())
		}
	}

	vs := v.Call(paramsV)

	if len(vs) == 0 {
		log.Warn("call method no return")
		return reflect.ValueOf(nil)
	}
	return vs[0].Interface()
}

func toPtr(obj interface{}) reflect.Value {
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	return vp
}

func (c *callerDefault) getFuncInfo(name string) *def.MethodInfo {
	if m, ok := def.MethodsPools[name]; ok {
		return &m
	}
	log.Errorf("not find name [%s]", name)
	return nil
}
