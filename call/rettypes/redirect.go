package rettypes

import (
	"bytes"
	"gitee.com/fast_api/api/def"
	"io"
	"net/http"
	"reflect"
)

type Redirect struct {
	location string
}

func NewRedirect(location string) *Redirect {
	return &Redirect{location: location}
}

func (r *Redirect) Code() int {
	return http.StatusFound
}

func (r *Redirect) Append(header def.ReadHeader) map[string]string {
	return map[string]string{"Location": r.location}
}

func (r *Redirect) ContentType() string {
	return def.CONTENT_HTML
}

func (r *Redirect) Return() io.Reader {
	return bytes.NewBufferString("")
}

func (r *Redirect) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*Redirect)(nil))}
}
