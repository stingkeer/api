package def

import (
	"net/http"
	"net/url"
)

type Request struct {
	*http.Request
	rw http.ResponseWriter
}

func (r *Request) ResponseWriter() http.ResponseWriter {
	return r.rw
}

func WithRequest(rw http.ResponseWriter, req *http.Request) *Request {
	return &Request{rw: rw, Request: req}
}

type Serialize interface {
	Encode(interface{}) *Content
	// Decode interface{} is out
	Decode([]byte, interface{}) error
}

type Match interface {
	Match(url *url.URL) *Entry
	Add(key string, data interface{})
}

type Caller interface {
	// CallerTrace
	// Call request ==> object
	Call(f *Entry, req *Request) interface{}
}

type CallerTrace interface {
	Before() bool
	After()
}
