package public

const (
	Content_JSON = "application/json"
	HEAD_CONST   = "API_HEADER_TYPE"
)

type Header interface {
	Add(key, value string)
	Get(key string) string
	Values(key string) []string
}
