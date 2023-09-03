package core

import (
	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/match"
	"gitee.com/fast_api/api/serialize"
)

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
	http.AddHttpHandle(http.NewApiIntercept(def.DefaultContext.Match, def.DefaultContext.Caller, def.DefaultContext.Serialize, pool))

	//
	http.AddHttpHandle(http.DefaultStatic)
}
