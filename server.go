package api

import (
	"net/http"
)

func Start(addr string) {
	initDef()
	apiServer := ApiService{
		&MatchImpl{GetApi().getMaps()},
		&JSONConvertImpl{},
		&CallerDefault{},
	}
	http.ListenAndServe(":8080", &apiServer)
}

type ApiService struct {
	match   Match
	convert Convert
	caller  Caller
}

func (a *ApiService) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fun, params := a.match.match(req.URL, req.Method)
	if fun != nil {
		inf := a.caller.call(fun, params)
		rw.Write(a.convert.convert(inf))
	}
}
