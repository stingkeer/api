package rettypes

import (
	"bytes"
	"io"
	"net/http"
	"reflect"

	"go.aew.app/api.v1/def"
)

var (
	_ def.AppendHeader = (*Redirect)(nil)
	_ def.RetAdapter   = (*Redirect)(nil)
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
