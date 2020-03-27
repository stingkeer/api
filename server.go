package api

import (
	"net/http"
)

func Start(addr string) {
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
	fun := a.match.match(req.URL, req.Method)
	if fun != nil {
		inf := a.caller.call(fun)
		rw.Write(a.convert.convert(inf))
	}
}

type Entry struct {
	url    string
	group  string
	method string
	f      interface{}
}
