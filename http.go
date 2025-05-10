//go:build !http3

package api

import (
	"crypto/tls"
	"net"
	"net/http"

	"go.aew.app/api/def"
	_ "go.aew.app/api/kit"
	"go.aew.app/api/log"
)

var server = NewServer(def.DefaultContext.Pool)

func StartService(ops ...Optional) error {
	apply(&defaultConf, ops...)
	log.Infof("listen addr %s", defaultConf.listen)
	server := &http.Server{
		Addr:    defaultConf.listen,
		Handler: server,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func StartTLSService(ops ...Optional) error {
	apply(&defaultConf, ops...)
	log.Infof("listen tls addr %s", defaultConf.listen)
	t, err := loadTls(&defaultConf)
	if err != nil {
		log.Error(err)
		return err
	}
	server := &http.Server{
		Addr:      defaultConf.listen,
		Handler:   server,
		TLSConfig: t,
	}
	config := t.Clone()

	// if !strSliceContains(config.NextProtos, "http/1.1") {
	// 	config.NextProtos = append(config.NextProtos, "http/1.1")
	// }

	// if server.shuttingDown() {
	// 	return http.ErrServerClosed
	// }

	addr := server.Addr
	if addr == "" {
		addr = ":https"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error(err)
		return err
	}
	defer ln.Close()
	tlsListener := tls.NewListener(ln, config)
	if err := server.Serve(tlsListener); err != nil {
		log.Error(err)
	}
	return nil
}
