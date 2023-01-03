package def

import "net/http"

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
