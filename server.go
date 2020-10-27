package api

import (
	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/public"
	"gitee.com/fast_api/api/transverter"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
	"sync"
)

func Start(addr string) {
	logrus.Infof("server listem %s", addr)
	PackApi()
	http.ListenAndServe(addr, DefaultService())
}

type Service struct {
	match   public.Match
	convert public.ResultConvert
	caller  public.Caller
}

func DefaultService() *Service {
	convert := &transverter.JSONConvertImpl{}
	caller := call.NewCaller(convert, &transverter.DefaultTypeConvert{})
	apiServer := &Service{
		&MatchImpl{GetApi().getStore()},
		convert,
		caller,
	}
	return apiServer
}

func (a *Service) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ApiHttp(rw, req, func() *Service {
		return a
	})
}

var one sync.Once
var a *Service

func ApiHttp(rw http.ResponseWriter, req *http.Request, service func() *Service) {
	one.Do(func() {
		if service == nil {
			a = DefaultService()
		} else {
			a = service()
		}
	})
	if a == nil {
		logrus.Error("service is nil")
	}
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
