package http

import (
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"
)

var DefaultStatic = NewStatic()

type (
	StaticOps interface {
		Rewrite(orig string, replace string)
	}
	staticEntry struct {
		fs      http.FileSystem
		dirPath string
		rewrite map[string]string
	}
	StaticOption func(StaticOps)
)

func StaticRewrite(orig string, replace string) StaticOption {
	return func(ops StaticOps) {
		ops.Rewrite(orig, replace)
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

func (s *Static) Http(rw http.ResponseWriter, req *http.Request) bool {
	for reg, handler := range s.m {
		if ok, _ := s.match(req.URL.Path, reg); ok {
			return handler.ServeHTTP(rw, req)
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
	s.m[path] = entry
	for _, op := range sOps {
		op(&entry)
	}
}

func (s *Static) Order() int {
	return 99
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
