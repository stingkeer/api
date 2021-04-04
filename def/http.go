package def

const (
	Content_JSON        = "application/json"
	CONTENT_STREAM      = "application/octet-stream"
	HEAD_CONST          = "API_HEADER_TYPE"
	CONTENT_DISPOSITION = "Content-Disposition"
)

//used in param
//eg.
/**

 GET(func(a def.Header) interface{} {
		return a.Values("Accept-Encoding")
 }, "/h")

*/
type Header interface {
	ReadHeader
	WriteHeader
}

//used in retType
type ContentType interface {
	Content() string
}

type ReadHeader interface {
	Get(key string) string
	Values(key string) []string
}

type WriteHeader interface {
	Add(key, value string)
}

//used in retType
type AppendHeader interface {
	Append(header ReadHeader) map[string]string
}

type HttpStatus interface {
	Code() int
}
