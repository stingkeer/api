package rettypes

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/log"
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
	fileSize    int64
	name        *string
}

// Read(p []byte) (n int, err error)
func (s *Stream) Read(p []byte) (n int, err error) {
	iTotal := s.end - s.start + 1
	if s.readTotal+int64(len(p)) < iTotal {
		t, e := s.rLimit.Read(p)
		s.readTotal += int64(t)
		return t, e
	} else {
		lRead, e := s.rLimit.Read(p)
		if e != nil {
			return lRead, e
		}
		return int(iTotal - s.readTotal), io.EOF
	}

}

func (s *Stream) Return() io.Reader {
	if s.end == 0 {
		return s.rLimit
	}
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

func (s *Stream) SetRateLimit(bytesPerSec float64) *Stream {
	s.rLimit.SetRateLimit(bytesPerSec)
	return s
}

func (s *Stream) ContentType() string {
	if s.contextType != nil {
		return *s.contextType
	}
	if seeker, sb := s.io.(io.Seeker); sb {
		if s.contextType == nil {
			temp := make([]byte, 512)
			if _, err := s.io.Read(temp); err != nil && err != io.EOF {
				panic(err)
			}
			cTyp := http.DetectContentType(temp)
			s.contextType = &cTyp
			seeker.Seek(0, io.SeekStart)
		}
		return *s.contextType
	}
	return def.CONTENT_STREAM
}

func (s *Stream) SetContentType(conTyp string) {
	s.contextType = &conTyp
}

func (s *Stream) Strings() string {
	return "Stream"
}

func (s *Stream) Close() error {
	if v, b := s.io.(io.Closer); b {
		log.Trace("close file ")
		return v.Close()
	}
	return nil
}

// bytes=0-1023
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
	if seeker, sb := s.io.(io.Seeker); sb && s.fileSize == 0 {

		size, err := seeker.Seek(0, io.SeekEnd)
		if err == nil {
			s.fileSize = size
			if s.end == 0 {
				s.end = s.fileSize - 1
			}
			seeker.Seek(0, io.SeekStart)

			//Accept-Ranges: bytes
			s.AddHeader("Accept-Ranges", "bytes")
		}
	}

	rge := header.Get("Range")
	if rge != "" {
		if seeker, sb := s.io.(io.Seeker); sb {
			s.parseRange(rge)
			//HTTP/1.1 206 Partial Content
			//Content-Range: bytes 0-1023/146515
			//Content-Length: 1024

			if s.end-s.start > s.fileSize {
				panic("more than size")
			}
			s.AddHeader("Content-Length", fmt.Sprintf("%d", s.end-s.start+1))
			s.AddHeader("Content-Range", fmt.Sprintf("bytes %d-%d/%d", s.start, s.end, s.fileSize))
			if s.start > 0 {
				seeker.Seek(s.start, io.SeekCurrent)
			}
			s.code = http.StatusPartialContent
			s.isRange = true
			return s.heads
		}
	} else {
		if s.fileSize != 0 {
			s.AddHeader("Content-Length", fmt.Sprintf("%d", s.fileSize))
		}
	}
	s.reSetFileName()
	return s.heads
}

func (s *Stream) Code() int {
	if s.code != 0 {
		return s.code
	}
	return http.StatusOK
}

func (s *Stream) SetCode(code int) {
	s.code = code
}

func (s *Stream) AddHeader(k, v string) {
	s.heads[k] = v
}

func (s *Stream) SetName(name string) *Stream {
	s.name = &name
	return s
}

func (s *Stream) reSetFileName() {
	if _, b := s.heads[def.CONTENT_DISPOSITION]; b {
		return
	}
	if s.name != nil {
		s.AddHeader(def.CONTENT_DISPOSITION, fmt.Sprintf("attachment; filename=%s", *s.name))
		return
	}
	if v, b := s.io.(*os.File); b && *s.contextType == def.CONTENT_STREAM {
		s.AddHeader(def.CONTENT_DISPOSITION, fmt.Sprintf("attachment; filename=%s", path.Base(v.Name())))
	}
}
