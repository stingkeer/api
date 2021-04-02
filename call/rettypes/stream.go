package rettypes

import (
	"bytes"
	"fmt"
	"gitee.com/fast_api/api/def"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"reflect"
)

type Stream struct {
	typ   int
	f     *os.File
	b     []byte
	heads map[string]string
}

func (s Stream) Return() io.Reader {
	switch s.typ {
	case 1:
		return s.f
	case 2:
		return bytes.NewBuffer(s.b)
	default:
		panic("stream not support")
	}
}

func (s Stream) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*Stream)(nil)).Elem()}
}

func NewFileStream(file *os.File) Stream {
	return Stream{f: file, typ: 1, heads: make(map[string]string)}
}

func NewBytesStream(b []byte) Stream {
	return Stream{b: b, typ: 2, heads: make(map[string]string)}
}

func (s Stream) limit(int2 int) {

}

func (s Stream) Content() string {
	return def.CONTENT_STREAM
}

func (s Stream) Strings() string {
	return "Stream"
}

func (s Stream) Close() {
	if s.typ == 1 {
		logrus.Tracef("close file %s", s.f.Name())
		s.f.Close()
	}
}

func (s Stream) Append() map[string]string {
	if _, b := s.heads[def.CONTENT_DISPOSITION]; !b && s.typ == 1 {
		s.AddHeader(def.CONTENT_DISPOSITION, fmt.Sprintf("attachment; filename=%s", path.Base(s.f.Name())))
	}
	return s.heads
}

func (s Stream) AddHeader(k, v string) {
	s.heads[k] = v
}
