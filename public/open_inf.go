package public

import (
	"net/http"
	"net/url"
	"reflect"
)

type TypeConvert interface {
	//convert string param to value
	ConvertTo(value string, typ reflect.Type) reflect.Value
}

type Serialize interface {
	Encode(interface{}) *Content
	//interface{} is out
	Decode([]byte, interface{}) error
}

type Match interface {
	Match(url *url.URL) *Entry
	Add(key string, data interface{})
}

type Caller interface {
	//request ==> object
	Call(f *Entry, req *http.Request) interface{}
}
