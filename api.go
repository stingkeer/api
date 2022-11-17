package api

import (
	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/call/rettypes"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/dwarf"
	"gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/match"
	"gitee.com/fast_api/api/mg"
	"gitee.com/fast_api/api/serialize"
	stdhttp "net/http"
	"reflect"
	"runtime"
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
		mg.Invoke(func(match def.Match) {
			match.Add(url, e)
		})
		mg.Invoke(func(pool *def.MethodsPools) {
			v := reflect.ValueOf(f)
			fName := runtime.FuncForPC(v.Pointer()).Name()
			if pool.Get(fName) == nil {
				mg.Invoke(func(dwarfMaker *dwarf.DwarfMaker) {
					findName := dwarfMaker.LookFun(f)
					var args = make(map[string]dwarf.ArgsMeta)
					for _, arg := range findName.Args {
						args[arg.Name] = arg
					}
					pool.Set(findName.MethodName, &def.MethodInfo{
						Method:     f,
						MethodName: findName.MethodName,
						Param:      args,
					})
				})
			}
		})
	}
}

func getFnCaches() []*def.Entry {
	return fnCaches
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

}
