package http

import "net/http"

type readHead struct {
	req *http.Request
}

func (r *readHead) Get(key string) string {
	return r.req.Header.Get(key)
}

func (r *readHead) Values(key string) []string {
	return r.req.Header.Values(key)
}
