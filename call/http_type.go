package call

import (
	"gitee.com/fast_api/api/public"
	"net/http"
	"reflect"
	"strings"
)

type HttpType struct {
}

func (f HttpType) Mapper(param public.ParamWarp) reflect.Value {
	return reflect.ValueOf(param.Request)
}

func (f HttpType) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*http.Request)(nil)).Elem()}
}

type headType struct {
	req *http.Request
}

func (f headType) Mapper(param public.ParamWarp) reflect.Value {
	f.req = &param.Request
	return reflect.ValueOf(f)
}

func (f headType) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*public.Header)(nil)).Elem()}
}

func (f headType) Add(key, value string) {
	var b strings.Builder
	orig := f.Get(public.HEAD_CONST)
	if orig != "" {
		b.WriteString(orig)
		b.WriteByte(',')
	}
	b.WriteString(key + "=" + value)
	f.req.Header.Add(public.HEAD_CONST, b.String())
}

func (f headType) Get(key string) string {
	return f.req.Header.Get(key)
}

func (f headType) Values(key string) []string {
	return f.req.Header.Values(key)
}
