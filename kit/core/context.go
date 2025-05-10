package core

import (
	_ "unsafe"

	"go.aew.app/api/call"
	"go.aew.app/api/def"
	"go.aew.app/api/http"
	"go.aew.app/api/intercept"
	"go.aew.app/api/kit/handler/sgzip"
	"go.aew.app/api/match"
	"go.aew.app/api/serialize"
)

//go:linkname addHttpHandle go.aew.app/api/http.addHttpHandle
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
