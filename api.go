package api

import (
	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/call/rettypes"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/match"
	"gitee.com/fast_api/api/serialize"
	"gitee.com/fast_api/api/server"
	stdhttp "net/http"
)

type (
	httpMethod func(f interface{}, url string)
)

var (
	GET  = httpM(stdhttp.MethodGet)
	POST = httpM(stdhttp.MethodPost)
	PUT  = httpM(stdhttp.MethodPut)

	// RegisterErrorHandler error handler
	RegisterErrorHandler = http.RegisterErrorHandler

	// AddHttpHandle http handler
	AddHttpHandle = http.AddHttpHandle

	// RegisterTypeMapper type handler
	RegisterTypeMapper = call.RegisterTypeMapper

	// RegisterReturnHandler register handler
	RegisterReturnHandler = http.RegisterReturnHandler

	NewStream = rettypes.NewStream
)

var fnCaches []*def.Entry

func httpM(method string) httpMethod {
	return func(f interface{}, url string) {
		e := &def.Entry{
			Url:    url,
			Method: method,
			Fn:     f,
			Ids:    make(map[string]string),
		}
		fnCaches = append(fnCaches, e)
		server.Invoke(func(f def.Match) {
			f.Add(url, e)
		})
	}
}

func getFnCaches() []*def.Entry {
	return fnCaches
}

func init() {

	//default
	server.Provide(func() def.Serialize {
		return &serialize.JsonConvertImpl{}
	})

	server.Provide(func(resultConvert def.Serialize) def.Caller {
		return call.NewCaller(resultConvert)
	})

	server.Provide(func() def.Match {
		return match.NewMatchImpl()
	})

	server.Invoke(func(match def.Match, caller def.Caller, serialize def.Serialize) {
		AddHttpHandle(http.NewApiIntercept(match, caller, serialize))
	})

}
