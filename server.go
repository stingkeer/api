package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func Start(addr string) {
	logrus.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)
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
		h := a.convert.convert(inf)
		rw.Header().Add("Content-Type", h.ContentType)
		rw.Write(h.bytes)
	} else {

	}
}
