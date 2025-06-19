package core

import (
	"os"
	"strings"
	"sync"

	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/dwarf"
	"go.aew.app/api.v1/log"
)

var (
	maker = dwarf.NewDwarfMaker()
	once  sync.Once
)

func printArgs(args []dwarf.ArgsMeta) string {
	var s strings.Builder
	l := len(args) - 1
	for i, arg := range args {
		s.WriteString(arg.Name)
		if i != l {
			s.WriteString(",")
		}
	}
	return s.String()
}

func makerInit() {
	if dll := os.Getenv("API_DLL"); dll == "" {
		maker.Init(nil)
	} else {
		maker.Init(&dll)
	}
}

func HttpM(method string, ctx *def.Context) def.HttpMethod {
	//init dwarf
	once.Do(makerInit)
	return func(f interface{}, url string) def.Option {
		entry := &def.Entry{
			Url:        url,
			HttpMethod: method,
			Fn:         f,
		}
		ctx.Match.Add(url, entry)
		findM, err := maker.LookFun(entry.Fn)
		if err != nil {
			panic(err)
		}
		var args = make(map[string]dwarf.ArgsMeta)
		for _, arg := range findM.Args {
			args[arg.Name] = arg
		}
		methodInfo := &def.MethodInfo{
			Method:     entry,
			MethodName: findM.MethodName,
			Param:      args,
		}
		ctx.Pool.Set(findM.MethodName, methodInfo)
		log.Infof("[%s] %s(%s) mapping url = %s", entry.HttpMethod, findM.MethodName, printArgs(findM.Args), entry.Url)
		op := option{url: url, method: method}
		op.SetMethod(methodInfo)
		op.SetContext(ctx)
		return &op
	}
}

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
