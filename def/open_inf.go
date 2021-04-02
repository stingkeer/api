package def

import (
	"net/http"
	"net/url"
)

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
