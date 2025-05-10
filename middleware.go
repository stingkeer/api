package api

import (
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/kit/core"
)

// def.MiddleWare is func(req *http.Request) (ret any)
// When ret is nil, continue to the next interceptor
// Otherwise, ret will be applied until exit

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

func Routes(os ...def.Option) *MiddlewareOps {
	return AddRoutes(os...)
}

func RoutesMiddleware(os ...def.Option) *MiddlewareOps {
	return AddRoutes(os...)
}
