package def

import (
	"net/http"
)

type (
	HttpMethod func(f any, url string) *Option
	MiddleWare func(rw http.ResponseWriter, req *http.Request) bool
)

var DefaultContext *Context

type Context struct {
	Match     Match
	Pool      *MethodsPools
	Caller    Caller
	Serialize Serialize
}

type Option struct {
	mi  *MethodInfo
	ctx *Context
}

func (o *Option) SetContext(ctx *Context) *Option {
	o.ctx = ctx
	return o
}

func (o *Option) SetMethod(md *MethodInfo) *Option {
	o.mi = md
	return o
}

func (o *Option) Swagger(opsFn func(swagger SwaggerOps)) *Option {
	opsFn(&swaggerImpl{o.mi})
	return o
}

func (o *Option) SetMiddleware(m ...MiddleWare) *Option {
	o.mi.Middleware = append(o.mi.Middleware, m...)
	return o
}

const (
	Content_JSON        = "application/json"
	CONTENT_STREAM      = "application/octet-stream"
	CONTENT_HTML        = "text/html"
	HEAD_CONST          = "_API_HEADER_TYPE_"
	CONTENT_DISPOSITION = "Content-Disposition"
)

type Cookie interface {
	SetCookie(cookie *http.Cookie)
	Cookie(name string) (*http.Cookie, error)
}

// Header used in param
//eg.
/**

 GET(func(a def.Header) interface{} {
		return a.Values("Accept-Encoding")
 }, "/h")

*/
type Header interface {
	Cookie
	ReadHeader
	WriteHeader
}

// ContentType used in retType
type ContentType interface {
	ContentType() string
}

type ReadHeader interface {
	Get(key string) string
	Values(key string) []string
}

type WriteHeader interface {
	Add(key, value string)
}

// AppendHeader used in retType
type AppendHeader interface {
	Append(header ReadHeader) map[string]string
}

type HttpStatus interface {
	Code() int
}
