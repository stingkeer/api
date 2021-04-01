package api

import (
	ihttp "gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/server"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

func StartService(addr string) {
	PackApi()
	server.Invoke(func(ser *Service) {
		logrus.Infof("listen addr %s", addr)
		logrus.Error(http.ListenAndServe(addr, ser))
	})
}

type Service struct{}

func (ad *Service) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ApiHttp(rw, req, func() *Service {
		return a
	})
}

func init() {
	server.Provide(func() *Service {
		a = &Service{}
		return a
	})

}

var (
	one sync.Once
	a   *Service
)

func ApiHttp(rw http.ResponseWriter, req *http.Request, service func() *Service) {
	one.Do(func() {
		if service != nil {
			a = service()
		}
	})
	ihttp.DoHttp(rw, req)
}
