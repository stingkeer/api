//go:build http3

package api

import (
	"gitee.com/fast_api/api/def"
	_ "gitee.com/fast_api/api/kit"
	"gitee.com/fast_api/api/log"
	"github.com/quic-go/quic-go/http3"
)

var server = NewServer(def.DefaultContext.Pool)

func StartService(ops ...Optional) {
	panic("http3 only support tls")
}

func StartTLSService(ops ...Optional) {
	apply(&defaultConf, ops...)
	log.Infof("listen addr %s")

	t, err := loadTls(&defaultConf)
	if err != nil {
		log.Error(err)
		return
	}
	server := &http3.Server{
		Addr:      defaultConf.listen,
		Handler:   server,
		TLSConfig: t,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error(err)
	}
}
