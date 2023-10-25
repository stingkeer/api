package api

import (
	"net/http"

	"gitee.com/fast_api/api/def"
	_ "gitee.com/fast_api/api/kit"
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

type MiddlewareOps struct {
	ops []*def.Option
}

func (m *MiddlewareOps) Middleware(mw ...def.MiddleWare) *MiddlewareOps {
	for i := 0; i < len(m.ops); i++ {
		m.ops[i].SetMiddleware(mw...)
	}
	return m
}

func (m *MiddlewareOps) WithPrefix(...def.MiddleWare) *MiddlewareOps {
	return m
}

func AddRoutes(os ...*def.Option) *MiddlewareOps {
	return &MiddlewareOps{ops: os}
}
