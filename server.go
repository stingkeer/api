package api

import (
	"crypto/tls"
	"crypto/x509"
	ihttp "gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/server"
	"github.com/sirupsen/logrus"
	"io/ioutil"
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

func StartTLSService(addr string, caFile, certFile, keyFile string) {
	PackApi()
	server.Invoke(func(ser *Service) {
		logrus.Infof("listen addr %s", addr)
		caCertPool := x509.NewCertPool()
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			panic(err)
		}
		caCertPool.AppendCertsFromPEM(caCert)
		server := &http.Server{Addr: addr, Handler: ser,
			TLSConfig: &tls.Config{
				ClientAuth: tls.RequireAndVerifyClientCert,
				ClientCAs:  caCertPool,
			},
		}
		logrus.Error(server.ListenAndServeTLS(certFile, keyFile))
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
