package core

import (
	_ "unsafe"

	"go.aew.app/api.v1/call"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/http"
	"go.aew.app/api.v1/intercept"
	"go.aew.app/api.v1/kit/handler/sgzip"
	"go.aew.app/api.v1/match"
	"go.aew.app/api.v1/serialize"
)

//go:linkname addHttpHandle go.aew.app/api.v1/http.addHttpHandle
func addHttpHandle(f intercept.HttpIntercept)

func init() {
	json := &serialize.JsonConvertImpl{}
	pool := &def.MethodsPools{}
	def.DefaultContext = &def.Context{
		Serialize: json,
		Match:     match.NewMatchImpl(),
		Pool:      pool,
		Caller:    call.NewWsCaller(json, pool),
	}
	//
	addHttpHandle(http.NewApiIntercept(def.DefaultContext.Match, def.DefaultContext.Caller, def.DefaultContext.Serialize, pool))

	addHttpHandle(http.NewApiRespose())
	//
	addHttpHandle(http.NewNotFind(def.DefaultContext.Serialize))
	//
	addHttpHandle(http.DefaultStatic)
	//
	addHttpHandle(&sgzip.GZip{})
}
