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
	Add(key, value string)
	Get(key string) string
	Values(key string) []string
}

//used in retType
type ContentType interface {
	Content() string
}

//used in retType
type AppendHeader interface {
	Append() map[string]string
}
