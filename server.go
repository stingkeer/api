package api

import (
	"gitee.com/fast_api/api/mg"
	"net/http"
)

var server *Server

func init() {
	mg.Invoke(func(s *Server) {
		server = s
	})
}

func StartService(f ConfigFun) {
	if f != nil {
		server.SetConfig(*f(server.Config()))
	}
	server.ListenAndServe()
}

func StartTLSService(f ConfigFun) {
	if f != nil {
		server.SetConfig(*f(server.Config()))
	}
	server.StartTLSService()
}

// Http this combined with other http handler
func Http(rw http.ResponseWriter, req *http.Request) {
	server.ApiHttp(rw, req)
}
