package api

import (
	"crypto/tls"
	"net/http"
	"os"

	"go.aew.app/api/def"
	"go.aew.app/api/dwarf"
	ihttp "go.aew.app/api/http"
	"go.aew.app/api/log"
)

type Optional func(conf *ServerConfig)

var (
	_           http.Handler = (*ServerMain)(nil)
	defaultConf              = ServerConfig{
		dwarf:  dwarf.NewDwarfMaker(),
		listen: "0.0.0.0:8080",
	}
)

func loadTls(s *ServerConfig) (tconf *tls.Config, err error) {
	if s.tlsConfig != nil {
		s.tlsConfig = &tls.Config{}
	}
	certs := make([]tls.Certificate, 1)
	certs[0], err = tls.X509KeyPair(defaultConf.certPEMBlock, defaultConf.keyPEMBlock)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	s.tlsConfig.Certificates = certs
	return s.tlsConfig, nil
}

type ServerMain struct {
	maker *dwarf.DwarfMaker
	pool  *def.MethodsPools
}

func (ad *ServerMain) Maker() *dwarf.DwarfMaker {
	return ad.maker
}

func NewServer(pool *def.MethodsPools) *ServerMain {
	return &ServerMain{pool: pool}
}

func (ad *ServerMain) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ihttp.DoHttp(rw, req)
}

func WithListen(s string) Optional {
	return func(conf *ServerConfig) {
		conf.listen = s
	}
}

func WithTLSConfig(tlsConfig *tls.Config) Optional {
	return func(conf *ServerConfig) {
		conf.tlsConfig = tlsConfig
	}
}

func WithTLS(certPEMBlock, keyPEMBlock []byte) Optional {
	return func(conf *ServerConfig) {
		conf.certPEMBlock = certPEMBlock
		conf.keyPEMBlock = keyPEMBlock
	}
}

func WithTLSFile(certFile, keyFile string) Optional {
	return func(conf *ServerConfig) {
		certPEMBlock, err := os.ReadFile(certFile)
		if err != nil {
			log.Error(err)
			return
		}
		keyPEMBlock, err := os.ReadFile(keyFile)
		if err != nil {
			log.Error(err)
			return
		}
		conf.certPEMBlock = certPEMBlock
		conf.keyPEMBlock = keyPEMBlock
	}
}

func WithCa(caPEMBlock []byte) Optional {
	return func(conf *ServerConfig) {
		conf.caPEMBlock = caPEMBlock
	}
}

func apply(s *ServerConfig, ops ...Optional) {
	if ops == nil {
		return
	}
	for i := 0; i < len(ops); i++ {
		if ops[i] != nil {
			ops[i](s)
		}

	}
}
