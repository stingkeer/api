package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func Start(addr string) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)
	initDef()
	j := &JSONConvertImpl{}
	apiServer := Service{
		&MatchImpl{GetApi().getMaps()},
		j,
		&CallerDefault{j},
	}
	logrus.Infof("server listem %s", addr)
	http.ListenAndServe(addr, &apiServer)
}

type Service struct {
	match   Match
	convert Convert
	caller  Caller
}

func (a *Service) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error(err)
		}
	}()
	logrus.Tracef("incoming req method [%s] , url [%s]", req.Method, req.URL.String())
	fun := a.match.match(req.URL)
	if fun != nil {
		inf := a.caller.call(fun, req)
		h := a.convert.convertTo(inf)
		rw.Header().Add("Content-Type", h.ContentType)
		rw.Write(h.bytes)
	}
}
