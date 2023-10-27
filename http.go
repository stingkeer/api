//go:build !http3

package api

import (
	"net/http"

	"gitee.com/fast_api/api/def"
	_ "gitee.com/fast_api/api/kit"
	"gitee.com/fast_api/api/log"
)

var server = NewServer(def.DefaultContext.Pool)

func StartService(ops ...Optional) {
	apply(&defaultConf, ops...)
	server := &http.Server{
		Addr:    defaultConf.listen,
		Handler: server,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error(err)
	}
}

func StartTLSService(ops ...Optional) {
	apply(&defaultConf, ops...)
	t, err := loadTls(&defaultConf)
	if err != nil {
		log.Error(err)
		return
	}
	server := &http.Server{
		Addr:      defaultConf.listen,
		Handler:   server,
		TLSConfig: t,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error(err)
	}
}
