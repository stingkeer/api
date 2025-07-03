package core

import "go.aew.app/api.v1/def"

var (
	_ def.Option = (*option)(nil)
)

type option struct {
	mi          *def.MethodInfo
	ctx         *def.Context
	url, method string
}

// SetKV implements def.Option.
func (o *option) StoreKV(key string, v any) {
	o.mi.KV.Store(key, v)
}

// Method implements def.Option.
func (o *option) Method() string {
	return o.method
}

// Path implements def.Option.
func (o *option) Path() string {
	return o.url
}

func (o *option) SetContext(ctx *def.Context) def.Option {
	o.ctx = ctx
	return o
}

func (o *option) SetMethod(md *def.MethodInfo) def.Option {
	o.mi = md
	return o
}

func (o *option) Swagger(opsFn func(swagger def.SwaggerOps)) def.Option {
	opsFn(&swaggerImpl{mi: o.mi, SwaggerSecurit: SwaggerSecurit{Ops: []def.Option{o}}})
	return o
}

func (o *option) SetMiddleware(m ...def.MiddleWare) def.Option {
	o.mi.Middleware = append(o.mi.Middleware, m...)
	return o
}
