package call

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
	"strings"

	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/intercept"
	"gitee.com/fast_api/api/log"
	"gitee.com/fast_api/api/utils"
)

var _ def.Caller = (*callerDefault)(nil)

type callerDefault struct {
	serialize  def.Serialize
	pool       *def.MethodsPools
	mIntercept intercept.MethodIntercept
}

var (
	adapters       = make(map[reflect.Type]def.Adapter)
	adapterGeneric = make(map[string]def.Adapter)
)

func NewCaller(serialize def.Serialize, pool *def.MethodsPools) *callerDefault {
	return &callerDefault{
		serialize:  serialize,
		mIntercept: NewUserProxyInvokeImpl(methodInvokes),
		pool:       pool,
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

func (c *callerDefault) decodeToStruct(ptyp reflect.Type, bytes []byte) reflect.Value {
	newT := reflect.New(ptyp)
	//decode data to struct
	err1 := c.serialize.Decode(bytes, newT.Interface())
	if err1 != nil {
		panic(err1)
	}
	return newT.Elem()
}

func (c *callerDefault) doSlice(ptyp reflect.Type, bytes []byte) reflect.Value {

	switch ptyp.Elem().Kind() {
	case reflect.Uint8:
		return reflect.ValueOf(bytes)
	default:
		return c.decodeToStruct(ptyp, bytes)
	}
}

func (c *callerDefault) doBody(pw *def.ParamWarp) reflect.Value {

	bytes, err := io.ReadAll(pw.Body)
	if err != nil {
		panic(err)
	}

	switch pw.PTyp.Kind() {
	case reflect.String:
		return reflect.ValueOf(string(bytes))
	case reflect.Slice:
		return c.doSlice(pw.PTyp, bytes)
	}

	return c.decodeToStruct(pw.PTyp, bytes)
}

func (c *callerDefault) doStructParam(params url.Values, pw *def.ParamWarp) reflect.Value {
	t := pw.PTyp
	v := reflect.New(pw.PTyp).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			v.Field(i).Set(c.doStructParam(params, &def.ParamWarp{Request: pw.Request, PTyp: field.Type}))
			continue
		}
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		tags := strings.Split(jsonTag, ",")
		firstName := tags[0]
		if t, b := adapters[field.Type]; b {
			if pv := params.Get(firstName); pv != "" {
				mV := t.Mapper(&def.ParamWarp{Request: pw.Request, PTyp: field.Type, PName: jsonTag, PValue: pv})
				v.Field(i).Set(mV)
			}
		}
	}
	return v
}

// Call request == call(def) => value
func (c *callerDefault) Call(f *def.Entry, req *def.Request) interface{} {
	v := reflect.ValueOf(f.Fn)
	m := c.pool.FuncInfo(f.Fn)
	if m == nil {
		log.Error("not find method in exe")
		os.Exit(2)
	}
	params := req.URL.Query()
	//
	f.Ids.Range(func(key, value string) {
		params.Add(key, value)
	})

	paramsV := make([]reflect.Value, len(m.Param))

	bodyCount := 0

	for pName, p := range m.Param {
		pw := &def.ParamWarp{Request: *req}
		pw.PTyp = v.Type().In(p.Order)
		pw.PName = pName

		if bodyCount >= 2 {
			panic(fmt.Errorf("there can only be one structure in a method body %s", f.HttpMethod))
		}

		// Process body
		// only one struct type in method !
		// only name's body struct can convert to data
		if pName == "body" {
			paramsV[p.Order] = c.doBody(pw)
			//inc ptr for count struct number
			bodyCount++
			continue
		}

		//First determine the basic type assignment
		if t, b := adapters[p.Typ]; b {
			if param, exist := params[pName]; exist {
				pw.PValue = param[0]
			}
			paramsV[p.Order] = t.Mapper(pw)
			continue
		}

		// adapterGeneric
		if pw.PTyp.Kind() == reflect.Struct {
			tName, _ := TypeInfo(pw.PTyp.String())
			if td, b1 := adapterGeneric[tName]; b1 {
				if param, exist := params[pName]; exist {
					pw.PValue = param[0]
				}
				paramsV[p.Order] = td.Mapper(pw)
				continue
			}
		}

		if pw.PTyp.Kind() == reflect.Struct {
			paramsV[p.Order] = c.doStructParam(params, pw)
			continue
		}

		//default value
		log.Warnf("[not support %s ] set default value", pw.PTyp)
		// fmt.Println(pw.PTyp.Kind(), reflect.TypeOf((*def.IntReq)(nil)).Elem())
		paramsV[p.Order] = utils.DefaultCallValue(pw.PTyp)
	}

	var vs []reflect.Value
	if c.mIntercept == nil {
		vs = v.Call(paramsV)
	} else {
		vs = c.mIntercept.Invoke(m, paramsV)
	}
	if len(vs) == 0 {
		log.Warn("call method no return")
		return nil
	}
	x := vs[0]
	return x.Interface()
}

func toPtr(obj interface{}) reflect.Value {
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	return vp
}

func (c *callerDefault) getFuncInfo(name string) *def.MethodInfo {
	mInfo := c.pool.Get(name)
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
