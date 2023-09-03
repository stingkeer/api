package api

import (
	"gitee.com/fast_api/api/def"
	_ "gitee.com/fast_api/api/kit"
	"net/http"
)

var server = NewServer(def.DefaultContext.Pool)

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
