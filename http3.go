//go:build http3

package api

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"

	"gitee.com/fast_api/api/def"
	_ "gitee.com/fast_api/api/kit"
	"gitee.com/fast_api/api/log"
	"github.com/quic-go/quic-go/http3"
)

var server = NewServer(def.DefaultContext.Pool)

func StartService(ops ...Optional) error {
	return errors.New("http3 only support tls")
}

func StartTLSService(ops ...Optional) error {
	apply(&defaultConf, ops...)
	log.Infof("QUIC listen addr %s", defaultConf.listen)

	t, err := loadTls(&defaultConf)
	if err != nil {
		log.Error(err)
		return err
	}
	quicServer := &http3.Server{
		Addr:      defaultConf.listen,
		Handler:   server,
		TLSConfig: t,
	}

	hErr := make(chan error)
	qErr := make(chan error)
	go func() {
		hErr <- httpTls(&defaultConf, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			quicServer.SetQuicHeaders(w.Header())
			server.ServeHTTP(w, r)
		}))
	}()
	go func() {
		qErr <- quicServer.Serve(nil)
	}()

	select {
	case err := <-hErr:
		quicServer.Close()
		return err
	case err := <-qErr:
		// Cannot close the HTTP server or wait for requests to complete properly :/
		return err
	}
}

func httpTls(s *ServerConfig, handler http.Handler) error {
	t, err := loadTls(s)
	if err != nil {
		log.Error(err)
		return err
	}
	server := &http.Server{
		Addr:      s.listen,
		Handler:   handler,
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
