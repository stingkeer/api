package core

import (
	_ "unsafe"

	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/intercept"
	"gitee.com/fast_api/api/match"
	"gitee.com/fast_api/api/serialize"
)

//go:linkname addHttpHandle gitee.com/fast_api/api/http.addHttpHandle
func addHttpHandle(f intercept.HttpIntercept)

func init() {
	json := &serialize.JsonConvertImpl{}
	pool := &def.MethodsPools{}
	def.DefaultContext = &def.Context{
		Serialize: json,
		Match:     match.NewMatchImpl(),
		Pool:      pool,
		Caller:    call.NewCaller(json, pool),
	}
	//
	addHttpHandle(http.NewApiIntercept(def.DefaultContext.Match, def.DefaultContext.Caller, def.DefaultContext.Serialize, pool))

	//
	addHttpHandle(http.NewNotFind(def.DefaultContext.Serialize))
	//
	addHttpHandle(http.DefaultStatic)
}
