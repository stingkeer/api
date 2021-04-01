package api

import (
	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/public"
	"gitee.com/fast_api/api/server"
	stdhttp "net/http"
)

type httpMethod func(f interface{}, url string)

var (
	GET  = httpM(stdhttp.MethodGet)
	POST = httpM(stdhttp.MethodPost)
	PUT  = httpM(stdhttp.MethodPut)

	//error handler
	RegisterErrorHandler = http.RegisterErrorHandler

	//http handler
	AddHttpHandle = http.AddHttpHandle

	//type handler
	RegisterTypeMapper = call.RegisterTypeMapper
)

var fnCaches []*public.Entry

func httpM(method string) httpMethod {
	return func(f interface{}, url string) {
		e := &public.Entry{
			Url:    url,
			Method: method,
			Fn:     f,
			Ids:    make(map[string]string),
		}
		fnCaches = append(fnCaches, e)
		server.Invoke(func(f public.Match) {
			f.Add(url, e)
		})
	}
}

func getFnCaches() []*public.Entry {
	return fnCaches
}
