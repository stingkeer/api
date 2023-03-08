package call

import (
	"fmt"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/intercept"
	"gitee.com/fast_api/api/log"
	"gitee.com/fast_api/api/mg"
	"gitee.com/fast_api/api/utils"
	"io"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

type callerDefault struct {
	serialize  def.Serialize
	mIntercept intercept.MethodIntercept
}

var (
	adapters       = make(map[reflect.Type]def.Adapter)
	adapterGeneric = make(map[string]def.Adapter)
	pool           *def.MethodsPools
	once           sync.Once
)

func NewCaller(serialize def.Serialize) *callerDefault {
	return &callerDefault{
		serialize:  serialize,
		mIntercept: &defaultProxyInvoke{},
	}
}

func RegisterTypeMapper(adapter def.Adapter) {
	if adapter != nil {
		for _, m := range adapter.Register() {
			adapters[m] = adapter
		}
	}
}

func RegisterGenericTypeMapper(adapter def.Adapter) {
	if adapter != nil {
		for _, m := range adapter.Register() {
			g, _ := TypeInfo(m.String())
			adapterGeneric[g] = adapter
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
	f.Ids.Range(func(key, value any) bool {
		params.Add(key.(string), value.(string))
		return true
	})

	paramsV := make([]reflect.Value, len(m.Param))
	//TODO
	for pName, p := range m.Param {
		pw := def.ParamWarp{Request: *req}
		pw.PTyp = v.Type().In(p.Order)
		pw.PName = pName
		if t, b := adapters[p.Typ]; b {
			if param, exist := params[pName]; exist {
				pw.PValue = param[0]
			}
			paramsV[p.Order] = t.Mapper(pw)
		} else if pw.PTyp.Kind() == reflect.Struct {
			tName, _ := TypeInfo(pw.PTyp.String())
			if td, b1 := adapterGeneric[tName]; b1 {
				if param, exist := params[pName]; exist {
					pw.PValue = param[0]
				}
				paramsV[p.Order] = td.Mapper(pw)
			}
		} else if pw.PTyp.Kind() == reflect.Struct && req.Method == http.MethodPost {
			newT := reflect.New(pw.PTyp)
			bytes, err := io.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			err1 := c.serialize.Decode(bytes, newT.Interface())
			if err1 != nil {
				panic(err1)
			}
			paramsV[p.Order] = newT.Elem()
		} else { //default value
			log.Tracef("not support %s set default value", pw.PTyp)
			fmt.Println(pw.PTyp.Kind(), reflect.TypeOf((*def.IntReq)(nil)).Elem())
			paramsV[p.Order] = utils.DefaultCallValue(pw.PTyp.Kind())
		}
	}
	var vs []reflect.Value
	if c.mIntercept == nil {
		vs = v.Call(paramsV)
	} else {
		vs = c.mIntercept.Invoke(v, m, paramsV)
	}
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
	once.Do(func() {
		err := mg.Invoke(func(poolMd *def.MethodsPools) {
			pool = poolMd
		})
		if err != nil {
			panic(err)
		}
	})
	mInfo := pool.Get(name)
	if mInfo != nil {
		return mInfo
	}
	log.Errorf("not find name [%s]", name)
	return nil
}

func TypeInfo(name string) (typ string, generic string) {
	i := strings.Index(name, "[")
	if i > 0 {
		return name[0:i], name[i : len(name)-1]
	} else {
		return name, ""
	}
}
