package types

import (
	"encoding/json"
	"net/http"
	"reflect"

	"go.aew.app/api/def"
)

var (
	_ def.Adapter = (*HttpType)(nil)
	_ def.Adapter = (*HeadType)(nil)
)

type HttpType struct{}

func (f HttpType) Mapper(param *def.ParamWarp) reflect.Value {
	if param.PTyp.Kind() == reflect.Struct {
		return reflect.ValueOf(param.Request.Request).Elem()
	}
	if param.PTyp.Kind() == reflect.Ptr {
		return reflect.ValueOf(param.Request.Request)
	}
	return reflect.Zero(param.PTyp)
}

func (f HttpType) Register() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*http.Request)(nil)).Elem(),
		reflect.TypeOf((*http.Request)(nil)),
	}
}

type HeadType struct {
	req *http.Request
}

func (f *HeadType) SetCookie(cookie *http.Cookie) {
	if v := cookie.String(); v != "" {
		f.Add("Set-Cookie", v)
	}
}

func (f *HeadType) Cookie(name string) (*http.Cookie, error) {
	return f.req.Cookie(name)
}

func (f HeadType) Mapper(param *def.ParamWarp) reflect.Value {
	f.req = param.Request.Request
	return reflect.ValueOf(&f)
}

func (f HeadType) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*def.Header)(nil)).Elem()}
}

func (f HeadType) Add(key, value string) {
	orig := f.Get(def.HEAD_CONST)
	m := make(map[string]string)
	if orig != "" {
		err := json.Unmarshal([]byte(orig), &m)
		if err != nil {
			return
		}
	}
	m[key] = value
	bytes, err := json.Marshal(m)
	if err != nil {
		return
	}
	f.req.Header.Add(def.HEAD_CONST, string(bytes))
}

func (f HeadType) Get(key string) string {
	return f.req.Header.Get(key)
}

func (f HeadType) Values(key string) []string {
	return f.req.Header.Values(key)
}
