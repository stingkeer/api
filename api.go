package api

import (
	"gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/public"
	"gitee.com/fast_api/api/server"
)

type httpMethod func(f interface{}, url string)

var (
	GET  = httpM(public.GET)
	POST = httpM(public.POST)
	PUT  = httpM(public.PUT)

	//error handler
	RegisterErrorHandler = http.RegisterErrorHandler
	//http handler
	AddHttpHandle = http.AddHttpHandle
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
