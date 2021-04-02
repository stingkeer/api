package def

const (
	Content_JSON        = "application/json"
	CONTENT_STREAM      = "application/octet-stream"
	HEAD_CONST          = "API_HEADER_TYPE"
	CONTENT_DISPOSITION = "Content-Disposition"
)

type Header interface {
	Add(key, value string)
	Get(key string) string
	Values(key string) []string
}

type ContentType interface {
	Content() string
}

type AppendHeader interface {
	Append() map[string]string
}
