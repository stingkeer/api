package api

import (
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/kit/core"
)

type MiddlewareOps struct {
	ops []def.Option
}

func (m *MiddlewareOps) Middleware(mw ...def.MiddleWare) *MiddlewareOps {
	for i := 0; i < len(m.ops); i++ {
		m.ops[i].SetMiddleware(mw...)
	}
	return m
}
func (o *MiddlewareOps) Swagger(opsFn func(swagger def.SwaggerSecurity)) *MiddlewareOps {
	opsFn(&core.SwaggerSecurit{Ops: o.ops})
	return o
}
func (m *MiddlewareOps) WithPrefix(...def.MiddleWare) *MiddlewareOps {
	return m
}

func AddRoutes(os ...def.Option) *MiddlewareOps {
	return &MiddlewareOps{ops: os}
}
