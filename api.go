package api

import (
	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/call/rettypes"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/dwarf"
	"gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/log"
	"gitee.com/fast_api/api/match"
	"gitee.com/fast_api/api/mg"
	"gitee.com/fast_api/api/serialize"
	stdhttp "net/http"
	"os"
	"sync"
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

var (
	maker = dwarf.NewDwarfMaker()
	once  sync.Once
)

func httpM(method string) httpMethod {
	//init dwarf
	once.Do(func() {
		if dll := os.Getenv("API_DLL"); dll == "" {
			maker.Init(nil)
		} else {
			maker.Init(&dll)
		}
	})
	return func(f interface{}, url string) {
		entry := &def.Entry{
			Url:    url,
			Method: method,
			Fn:     f,
		}
		initFnCache.Add(entry)
		err := mg.Invoke(func(match def.Match) {
			match.Add(url, entry)
		})
		if err != nil {
			panic(err)
		}
		findM, err := maker.LookFun(entry.Fn)
		if err != nil {
			panic(err)
		}
		var args = make(map[string]dwarf.ArgsMeta)
		for _, arg := range findM.Args {
			args[arg.Name] = arg
		}
		err = mg.Invoke(func(pool *def.MethodsPools) {
			pool.Set(findM.MethodName, &def.MethodInfo{
				Method:     entry,
				MethodName: findM.MethodName,
				Param:      args,
			})
		})
		if err != nil {
			panic(err)
		}
		log.Infof("[%s] %s(%s) mapping url = %s", entry.Method, findM.MethodName, printArgs(findM.Args), entry.Url)
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
