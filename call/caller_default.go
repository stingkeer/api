package call

import (
	"gitee.com/fast_api/api/public"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
)

type callerDefault struct {
	serialize  public.Serialize
	typConvert public.TypeConvert
}

func NewCaller(serialize public.Serialize, typConvert public.TypeConvert) *callerDefault {
	return &callerDefault{
		serialize:  serialize,
		typConvert: typConvert,
	}
}

func (c *callerDefault) Call(f *public.Entry, req *http.Request) interface{} {
	switch req.Method {
	case public.GET:
		return c.callGet(f, req)
	case public.POST:
		return c.callPost(f, req)
	default:
		logrus.Warnf("method not support %s", req.Method)
	}
	return nil
}

func IsMultipart(t reflect.Type) bool {
	return t == reflect.TypeOf((*multipart.Reader)(nil)).Elem()
}

func (c *callerDefault) callPost(f *public.Entry, req *http.Request) interface{} {
	v := reflect.ValueOf(f.Fn)
	if v.Type().NumIn() > 1 {
		logrus.Error("not support param > 1")
		return nil
	}
	p0 := v.Type().In(0)
	newT := reflect.New(p0)
	if IsMultipart(p0) { //file
		reader, err := req.MultipartReader()
		if err != nil {
			logrus.Error(err)
		}
		newT.Elem().Set(reflect.ValueOf(*reader))
	} else { //small body
		bytes, _ := ioutil.ReadAll(req.Body)
		//instant json
		err := c.serialize.Decode(bytes, newT.Interface())
		if err != nil {
			panic(err)
		}
	}
	vs := v.Call([]reflect.Value{newT.Elem()})
	if len(vs) == 0 {
		logrus.Warn("call method no return")
		return nil
	}
	return vs[0].Interface()
}

func (c *callerDefault) callGet(f *public.Entry, req *http.Request) interface{} {
	name := runtime.FuncForPC(reflect.ValueOf(f.Fn).Pointer()).Name()
	logrus.Tracef("call function name [%s]", name)
	m := c.getFuncInfo(name)
	if m == nil {
		logrus.Error("not find method in header")
		os.Exit(2)
	}
	t := reflect.TypeOf(f.Fn)
	var pvs = make([]reflect.Value, t.NumIn())
	logrus.Tracef("method has param [%d]", t.NumIn())
	params := req.URL.Query()
	for k, v := range f.Ids {
		params.Add(k, v)
	}

	//convert param value to reflect.value
	for name, p := range m.Param {
		if v, b := params[name]; b {
			pvs[p.Order] = c.paramTypeConvert(v[0], t.In(p.Order))
		} else {
			pvs[p.Order] = c.defaultCallValue(t.In(p.Order).Kind())
		}
	}
	vs := reflect.ValueOf(f.Fn).Call(pvs)
	if len(vs) == 0 {
		logrus.Warn("call method no return")
		return nil
	}
	return vs[0].Interface()
}

func (c *callerDefault) paramTypeConvert(value string, dest reflect.Type) reflect.Value {
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
	case reflect.Ptr:
		typeConvert := c.paramTypeConvert(value, dest.Elem())
		if typeConvert.Type().Kind() == reflect.Ptr {
			logrus.Error("your convert return ptr")
		} else {
			return toPtr(typeConvert.Interface())
		}
	case reflect.Struct:
		gValue := c.typConvert.ConvertTo(value, dest)
		logrus.Debugf("convert type %s dest type %s", gValue.Type(), dest)
		return gValue
	default:
		logrus.Errorf("not find type %s", dest)

	}
	return reflect.ValueOf(nil)
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
