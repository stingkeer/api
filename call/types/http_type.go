package types

import (
	"gitee.com/fast_api/api/def"
	"net/http"
	"reflect"
	"strings"
)

type HttpType struct {
}

func (f HttpType) Mapper(param def.ParamWarp) reflect.Value {
	return reflect.ValueOf(param.Request)
}

func (f HttpType) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*http.Request)(nil)).Elem()}
}

type HeadType struct {
	req *http.Request
}

func (f HeadType) Mapper(param def.ParamWarp) reflect.Value {
	f.req = &param.Request
	return reflect.ValueOf(f)
}

func (f HeadType) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*def.Header)(nil)).Elem()}
}

func (f HeadType) Add(key, value string) {
	var b strings.Builder
	orig := f.Get(def.HEAD_CONST)
	if orig != "" {
		b.WriteString(orig)
		b.WriteByte(',')
	}
	b.WriteString(key + "=" + value)
	f.req.Header.Add(def.HEAD_CONST, b.String())
}

func (f HeadType) Get(key string) string {
	return f.req.Header.Get(key)
}

func (f HeadType) Values(key string) []string {
	return f.req.Header.Values(key)
}
