package types

import (
	"encoding/json"
	"gitee.com/fast_api/api/def"
	"net/http"
	"reflect"
)

type HttpType struct{}

func (f HttpType) Mapper(param def.ParamWarp) reflect.Value {
	return reflect.ValueOf(param.Request)
}

func (f HttpType) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*http.Request)(nil)).Elem()}
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

func (f HeadType) Mapper(param def.ParamWarp) reflect.Value {
	f.req = &param.Request
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
