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
	j := &JSONConvertImpl{}
	apiServer := ApiService{
		&MatchImpl{GetApi().getMaps()},
		j,
		&CallerDefault{j},
	}
	http.ListenAndServe(":8080", &apiServer)
}

type ApiService struct {
	match   Match
	convert Convert
	caller  Caller
}

func (a *ApiService) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logrus.Tracef("incoming req method [%s] , url [%s]", req.Method, req.URL.String())
	fun := a.match.match(req.URL)
	if fun != nil {
		inf := a.caller.call(fun, req)
		h := a.convert.convertTo(inf)
		rw.Header().Add("Content-Type", h.ContentType)
		rw.Write(h.bytes)
	} else {

	}
}
