package def

import (
	"net/http"
)

type (
	HttpMethod func(f any, url string) Option
	// MiddleWare
	// If ret == nil The next MiddleWare will continue
	// if ret != nil Ret will be used as the result
	MiddleWare func(req *http.Request) (ret any)
)

var DefaultContext *Context

type Context struct {
	Match     Match
	Pool      *MethodsPools
	Caller    Caller
	Serialize Serialize
}

type Option interface {
	SetContext(ctx *Context) Option
	SetMethod(md *MethodInfo) Option
	StoreKV(key string, v any)
	Swagger(opsFn func(swagger SwaggerOps)) Option
	SetMiddleware(m ...MiddleWare) Option
	Path() string
	Method() string
}

const (
	Content_Type        = "Content-Type"
	Content_Encoding    = "Content-Encoding"
	Accept_Encoding     = "Accept-Encoding"
	Content_JSON        = "application/json;charset=utf-8"
	CONTENT_STREAM      = "application/octet-stream"
	CONTENT_HTML        = "text/html;charset=utf-8"
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
