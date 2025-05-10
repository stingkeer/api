package rettypes

import (
	"bytes"
	"io"
	"net/http"
	"reflect"

	"go.aew.app/api.v1/def"
)

var (
	_ def.RetAdapter   = (*Resp)(nil)
	_ def.AppendHeader = (*Resp)(nil)
	_ def.HttpStatus   = (*Resp)(nil)
)

type Resp struct {
	h           map[string]string
	res         any
	code        int
	contentType string
	serialize   def.Serialize
	reader      io.Reader
}

// NewResp By default, use json ContextType
// The incoming resp will be serialized
func NewResp(resp any) *Resp {
	return &Resp{
		res:         resp,
		code:        http.StatusOK,
		contentType: def.Content_JSON,
		serialize:   def.DefaultContext.Serialize,
	}
}

func (r *Resp) SetCode(code int) *Resp {
	r.code = code
	return r
}

func (r *Resp) SetReader(reader io.Reader) *Resp {
	r.reader = reader
	return r
}

func (r *Resp) SetHeader(header map[string]string) *Resp {
	r.h = header
	return r
}

func (r *Resp) SetContentType(contentType string) *Resp {
	r.contentType = contentType
	return r
}

func (r *Resp) SetSerialize(serialize def.Serialize) *Resp {
	r.serialize = serialize
	return r
}

// Code implements def.HttpStatus.
// ReadOnly
func (r *Resp) Code() int {
	return r.code
}

// Append implements def.AppendHeader.
// ReadOnly
func (r *Resp) Append(header def.ReadHeader) map[string]string {
	return r.h
}

// ContentType implements def.RetAdapter.
// ReadOnly
func (r *Resp) ContentType() string {
	return r.contentType
}

// Register implements def.RetAdapter.
// ReadOnly
func (r *Resp) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*Resp)(nil))}
}

// Return implements def.RetAdapter.
// ReadOnly
func (r *Resp) Return() io.Reader {
	if r.reader == nil {
		context := r.serialize.Encode(r.res)
		return bytes.NewBuffer(context.Bytes)
	}
	return r.reader
}
