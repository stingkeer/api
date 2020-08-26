package api

import (
	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/public"
	"gitee.com/fast_api/api/transverter"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"runtime/debug"
)

func Start(addr string) {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)
	//logrus.SetReportCaller(true)
	initDef()
	convert := &transverter.JSONConvertImpl{}
	caller := call.NewCaller(convert, &transverter.DefaultTypeConvert{})
	apiServer := Service{
		&MatchImpl{GetApi().getStore()},
		convert,
		caller,
	}
	logrus.Infof("server listem %s", addr)
	http.ListenAndServe(addr, &apiServer)
}

type Service struct {
	match   public.Match
	convert public.ResultConvert
	caller  public.Caller
}

func (a *Service) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			logrus.Error(err)
		}
	}()
	logrus.Tracef("incoming req Method [%s] , Url [%s]", req.Method, req.URL.String())
	entry := a.match.Match(req.URL)
	if nil == entry {
		return
	}
	if req.Method != entry.Method {
		logrus.Warnf("not support Method %s", req.Method)
		return
	}
	if entry.Fn != nil {
		inf := a.caller.Call(entry, req)
		h := a.convert.ConvertTo(inf)
		rw.Header().Add("Content-Type", h.ContentType)
		rw.Write(h.Bytes)
	}
}
