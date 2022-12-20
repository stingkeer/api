package http

import (
	"net/http"
	"strings"
)

var DefaultStatic = NewStatic()

type Static struct {
	m map[string]http.Handler
}

func NewStatic() *Static {
	return &Static{m: make(map[string]http.Handler)}
}

func (s *Static) Http(rw http.ResponseWriter, req *http.Request) bool {
	for reg, handler := range s.m {
		if ok, path := s.match(req.URL.Path, reg); ok {
			req.URL.Path = path
			handler.ServeHTTP(rw, req)
			return true
		}
	}
	return false
}

func (s *Static) AddStatic(path string, fileSystem http.FileSystem) {
	s.m[path] = http.FileServer(fileSystem)
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
