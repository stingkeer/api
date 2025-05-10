package http

import (
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/intercept"
)

var (
	DefaultStatic           = NewStatic()
	_             StaticOps = (*staticEntry)(nil)
)

type (
	StaticOps interface {
		Rewrite(orig string, replace string)
		DefaultFile(file string)
	}
	staticEntry struct {
		fs          http.FileSystem
		dirPath     string
		rewrite     map[string]string
		defaultFile string
	}
	StaticOption func(StaticOps)
)

// DefaultFile implements StaticOps.
func (s *staticEntry) DefaultFile(file string) {
	s.defaultFile = file
}

func StaticRewrite(orig string, replace string) StaticOption {
	return func(ops StaticOps) {
		ops.Rewrite(orig, replace)
	}
}

func StaticDefaultFile(file string) StaticOption {
	return func(ops StaticOps) {
		ops.DefaultFile(file)
	}
}

func (s *staticEntry) Rewrite(orig string, replace string) {
	s.rewrite[orig] = replace
}

func (s *staticEntry) ServeHTTP(rw http.ResponseWriter, req *http.Request) bool {
	if req.URL.Path == "/" {
		req.URL.Path = "/index.html"
	}
	//rewrite path
	pathWriter := req.URL.Path
	for orig, s3 := range s.rewrite {
		compile, err := regexp.Compile(orig)
		if err != nil {
			panic(err)
		}
		if o1 := compile.FindString(pathWriter); o1 != "" {
			pathWriter = strings.Replace(pathWriter, o1, s3, 1)
		}
	}

	f, err := s.fs.Open(path.Join(s.dirPath, pathWriter))
	if err != nil {
		if s.defaultFile == "" {
			return false
		} else {
			f, err = s.fs.Open(path.Join(s.dirPath, s.defaultFile))
		}
	}
	if err != nil {
		return false
	}
	defer f.Close()
	d, err := f.Stat()
	if err != nil {
		return false
	}
	if d.IsDir() {
		return false
	}
	http.ServeContent(rw, req, d.Name(), time.Now(), f)
	return true
}

type Static struct {
	m map[string]staticEntry
}

func NewStatic() *Static {
	return &Static{m: make(map[string]staticEntry)}
}

func (s *Static) Http(rw http.ResponseWriter, req *http.Request, ctx *intercept.HttpContext) bool {
	for reg, handler := range s.m {
		if ok, _ := s.match(req.URL.Path, reg); ok {
			b := handler.ServeHTTP(rw, req)
			if b {
				ctx.SkipResponse()
			}
			return b
		}
	}
	return false
}

// HandleStatic
// path is the mapping url
// dirPath is the real path
// fileSystem open[ join(path + dirPath) ]
func (s *Static) HandleStatic(path, dirPath string, fileSystem http.FileSystem, sOps ...StaticOption) {
	entry := staticEntry{
		fs:      fileSystem,
		dirPath: dirPath,
		rewrite: make(map[string]string),
	}

	for _, op := range sOps {
		op(&entry)
	}

	s.m[path] = entry
}

func (s *Static) Order() def.HandlerOrder {
	return def.Handler_STATIC
}

// /web/a.js  /web/*
func (s *Static) match(path string, reg string) (bool, string) {
	regs := strings.Split(reg, "/")
	paths := strings.Split(path, "/")
	for i, r := range regs {
		if r == "*" {
			return true, strings.Join(paths[i:], "/")
		}
		if paths[i] != r {
			return false, ""
		}
	}
	return true, strings.Join(paths, "/")
}
