package rettypes

import (
	"fmt"
	"gitee.com/fast_api/api/def"
	"github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

type Stream struct {
	io          io.Reader
	heads       map[string]string
	end         int64
	start       int64
	code        int
	readTotal   int64
	isRange     bool
	rLimit      *Reader
	contextType *string
}

//Read(p []byte) (n int, err error)
func (s *Stream) Read(p []byte) (n int, err error) {
	iTotal := s.end - s.start + 1
	if s.readTotal+int64(len(p)) < iTotal {
		t, e := s.rLimit.Read(p)
		s.readTotal += int64(t)
		return t, e
	} else {
		_, e := s.rLimit.Read(p)
		if e != nil {
			return int(iTotal - s.readTotal), e
		}
		return int(iTotal - s.readTotal), io.EOF
	}

}

func (s *Stream) Return() io.Reader {
	return s
}

func (s *Stream) Register() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*Stream)(nil))}
}

func NewStream(io io.Reader) *Stream {
	return &Stream{
		io:     io,
		heads:  make(map[string]string),
		rLimit: NewReader(io),
	}
}

func (s *Stream) SetRateLimit(bytesPerSec float64) {
	s.rLimit.SetRateLimit(bytesPerSec)
}

func (s *Stream) ContentType() string {
	if s.contextType == nil {
		st := def.CONTENT_STREAM
		s.contextType = &st
	}
	return *s.contextType
}

func (s *Stream) SetContentType(conTyp string) {
	s.contextType = &conTyp
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
func (s *Stream) parseRange(hRange string) {
	kv := strings.Split(hRange, "=")
	sEnd := strings.Split(kv[1], "-")
	if len(sEnd) < 1 {
		return
	}
	start, e := strconv.Atoi(sEnd[0])
	if e == nil {
		s.start = int64(start)
	}
	if sEnd[1] == "" {
		return
	}
	end, e := strconv.Atoi(sEnd[1])
	if e == nil {
		s.end = int64(end)
	}
}

func (s *Stream) Append(header def.ReadHeader) map[string]string {
	rge := header.Get("Range")
	if rge != "" {
		seeker, sb := s.io.(io.Seeker)
		if file, b := s.io.(fs.File); b && sb {
			s.parseRange(rge)
			//HTTP/1.1 206 Partial Content
			//Content-Range: bytes 0-1023/146515
			//Content-Length: 1024
			fstat, _ := file.Stat()
			if s.end == 0 {
				s.end = fstat.Size() - 1
			}
			logrus.Debug(s.start, s.end)
			//2-0
			if s.end-s.start > fstat.Size() {
				panic("more than size")
			}
			s.AddHeader("Content-Length", fmt.Sprintf("%d", s.end-s.start+1))
			s.AddHeader("Content-Range", fmt.Sprintf("bytes %d-%d/%d", s.start, s.end, fstat.Size()))
			if s.start > 0 {
				seeker.Seek(s.start, io.SeekCurrent)
			}
			s.code = http.StatusPartialContent
			s.isRange = true
			return s.heads
		}
	}
	if f, b := s.io.(fs.File); b {
		//Accept-Ranges: bytes
		stat, _ := f.Stat()
		s.AddHeader("Accept-Ranges", "bytes")
		s.AddHeader("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}
	if _, b := s.heads[def.CONTENT_DISPOSITION]; b {
		return s.heads
	}
	if v, b := s.io.(*os.File); b && *s.contextType == def.CONTENT_STREAM {
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
