package rettypes

import (
	"fmt"
	"gitee.com/fast_api/api/def"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

type Stream struct {
	io      io.Reader
	heads   map[string]string
	end     int64
	start   int64
	code    int
	total   int64
	isRange bool
	rLimit  *Reader
}

func (s *Stream) Read(p []byte) (n int, err error) {
	if s.total+int64(len(p)) < (s.end - s.start) {
		t, e := s.rLimit.Read(p)
		s.total += int64(t)
		return t, e
	} else {
		s.io.Read(p)
		return int((s.end - s.start) - s.total), io.EOF
	}

}

func (s *Stream) Return() io.Reader {

	if s.isRange {
		return s
	} else {
		return s.rLimit
	}
}

func (s *Stream) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*Stream)(nil))}
}

func NewStream(io io.Reader) *Stream {
	return &Stream{io: io, heads: make(map[string]string), rLimit: NewReader(io)}
}

func (s *Stream) SetRateLimit(bytesPerSec float64) {
	s.rLimit.SetRateLimit(bytesPerSec)
}

func (s *Stream) Content() string {
	return def.CONTENT_STREAM
}

func (s *Stream) Strings() string {
	return "Stream"
}

func (s *Stream) Close() error {
	if v, b := s.io.(io.Closer); b {
		logrus.Trace("close file ")
		return v.Close()
	}
	return nil
}

//bytes=0-1023
func getRange(s string) (int64, int64) {
	kv := strings.Split(s, "=")
	sEnd := strings.Split(kv[1], "-")
	start, _ := strconv.Atoi(sEnd[0])
	end, _ := strconv.Atoi(sEnd[1])
	return int64(start), int64(end)
}

func (s *Stream) Append(header def.ReadHeader) map[string]string {
	rge := header.Get("Range")
	if rge != "" {
		if file, b := s.io.(*os.File); b {
			s.start, s.end = getRange(rge)
			fmt.Println(s.start, s.end)
			//HTTP/1.1 206 Partial Content
			//Content-Range: bytes 0-1023/146515
			//Content-Length: 1024
			fstat, _ := file.Stat()
			if s.end-s.start > fstat.Size() {
				panic("more than size")
			}
			s.AddHeader("Content-Length", fmt.Sprintf("%d", s.end-s.start))
			s.AddHeader("Content-Range", fmt.Sprintf("bytes %d-%d/%d", s.start, s.end, fstat.Size()))
			file.Seek(s.start, io.SeekCurrent)
			s.code = http.StatusPartialContent
			s.isRange = true
			return s.heads
		}
	}
	if f, b := s.io.(*os.File); b {
		//Accept-Ranges: bytes
		stat, _ := f.Stat()
		s.AddHeader("Accept-Ranges", "bytes")
		s.AddHeader("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}
	if _, b := s.heads[def.CONTENT_DISPOSITION]; b {
		return s.heads
	}
	if v, b := s.io.(*os.File); b {
		s.AddHeader(def.CONTENT_DISPOSITION, fmt.Sprintf("attachment; filename=%s", path.Base(v.Name())))
	}
	return s.heads
}

func (s *Stream) Code() int {
	if s.code != 0 {
		return s.code
	}
	return http.StatusOK
}

func (s *Stream) AddHeader(k, v string) {
	s.heads[k] = v
}
