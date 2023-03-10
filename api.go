package api

import (
	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/call/rettypes"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/match"
	"gitee.com/fast_api/api/mg"
	"gitee.com/fast_api/api/serialize"
	stdhttp "net/http"
)

type (
	httpMethod func(f interface{}, url string)
)

const eg = `
please use api.GET or api.POST in init() method !
eg.
func init() {
	api.GET(func(username string) {
		fmt.Println(username)
	},"send")
}
`

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

	Html     = rettypes.NewHtml
	HtmlView = rettypes.HtmlView

	// Static static web
	Static = http.DefaultStatic.HandleStatic

	NewRedirect = rettypes.NewRedirect
)

func httpM(method string) httpMethod {
	if initFnCache.Init() {
		panic(eg)
	}
	return func(f interface{}, url string) {
		entry := &def.Entry{
			Url:    url,
			Method: method,
			Fn:     f,
		}
		initFnCache.Add(entry)
		mg.Invoke(func(match def.Match) {
			match.Add(url, entry)
		})
	}
}

func init() {

	//default
	mg.Provide(func() def.Serialize {
		return &serialize.JsonConvertImpl{}
	})

	mg.Provide(func(resultConvert def.Serialize) def.Caller {
		return call.NewCaller(resultConvert)
	})

	mg.Provide(func() def.Match {
		return match.NewMatchImpl()
	})

	mg.Invoke(func(match def.Match, caller def.Caller, serialize def.Serialize) {
		AddHttpHandle(http.NewApiIntercept(match, caller, serialize))
	})

	mg.Invoke(func(match def.Match, caller def.Caller, serialize def.Serialize) {
		AddHttpHandle(http.DefaultStatic)
	})
}
