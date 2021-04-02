package rettypes

import (
	"fmt"
	"gitee.com/fast_api/api/def"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"reflect"
)

type Stream struct {
	io    io.Reader
	heads map[string]string
}

func (s Stream) Return() io.Reader {
	return s.io
}

func (s Stream) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*Stream)(nil)).Elem()}
}

func NewStream(io io.Reader) Stream {
	return Stream{io: io, heads: make(map[string]string)}
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
	if v, b := s.io.(*os.File); b {
		logrus.Tracef("close file %s", v.Name())
		v.Close()
	}
}

func (s Stream) Append() map[string]string {
	if _, b := s.heads[def.CONTENT_DISPOSITION]; b {
		return s.heads
	}
	if v, b := s.io.(*os.File); b {
		s.AddHeader(def.CONTENT_DISPOSITION, fmt.Sprintf("attachment; filename=%s", path.Base(v.Name())))
	}
	return s.heads
}

func (s Stream) AddHeader(k, v string) {
	s.heads[k] = v
}
