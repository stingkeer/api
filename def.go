package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/dwarf"
	ihttp "gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/log"
	"gitee.com/fast_api/api/mg"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

type ConfigFun func(conf *Config) *Config

type Server struct {
	conf          Config
	prefix        string
	maker         *dwarf.DwarfMaker
	pool          *def.MethodsPools
	hasInitPacked bool
}

func (ad *Server) Maker() *dwarf.DwarfMaker {
	return ad.maker
}

func (ad *Server) Config() *Config {
	return &ad.conf
}

func init() {
	err := mg.Provide(NewServer)
	if err != nil {
		panic(err)
	}
}

func NewServer(pool *def.MethodsPools) *Server {
	config := Config{
		dwarf:  dwarf.NewDwarfMaker(),
		listen: ":8080",
	}
	return &Server{pool: pool, conf: config, maker: config.dwarf}
}

func (ad *Server) SetConfig(conf Config) {
	ad.conf = conf
	ad.maker = conf.dwarf
}

func (ad *Server) init(path *string) {
	usedMode := ad.maker.UsedMode()
	if usedMode.Mode() == dwarf.IncludeMode {
		if info, b := debug.ReadBuildInfo(); b {
			ad.maker.AddIncludeRegex(fmt.Sprintf("^%s/.+$", info.Path))
		} else {
			panic("It is forbidden to use `-s -w` during build")
		}
	}
	if dll, e := os.Executable(); e != nil && path == nil {
		ad.maker.Init(&dll)
	} else {
		ad.maker.Init(path)
	}
}

func (ad *Server) trimPrefix(s string) string {
	if s != "" {
		return strings.ReplaceAll(s, ad.prefix, "")
	}
	return s
}

func (ad *Server) ListenAndServe() {
	if !ad.hasInitPacked {
		ad.packInitApiWithPath(nil)
	}
	log.Infof("listen addr %s", ad.conf.Listen())
	log.Error(http.ListenAndServe(ad.conf.Listen(), ad))
}

func (ad *Server) StartTLSService() {
	if !ad.hasInitPacked {
		ad.packInitApiWithPath(nil)
	}
	log.Infof("listen addr %s", ad.conf.Listen())
	caCertPool := x509.NewCertPool()
	caCert, err := os.ReadFile(ad.conf.CaFile())
	if err != nil {
		panic(err)
	}
	caCertPool.AppendCertsFromPEM(caCert)
	serverListen := &http.Server{Addr: ad.conf.Listen(), Handler: ad,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  caCertPool,
		},
	}
	log.Error(serverListen.ListenAndServeTLS(ad.conf.CertFile(), ad.conf.KeyFile()))

}

func (ad *Server) ApiHttp(rw http.ResponseWriter, req *http.Request) {
	if !ad.hasInitPacked {
		ad.packInitApiWithPath(nil)
	}
	ihttp.DoHttp(rw, req)
}

func (ad *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ad.ApiHttp(rw, req)
}

func (ad *Server) SetLogTrimPrefix(prefixM string) {
	ad.prefix = prefixM
}

func (ad *Server) packInitApiWithPath(path *string) {
	start := time.Now()
	ad.init(path)
	log.Debugf("api had caches %d", initFnCache.Len())
	initFnCache.Range(func(index int, en *def.Entry) {
		findM, err := ad.maker.LookFun(en.Fn)
		if err != nil {
			panic(err)
		}
		var args = make(map[string]dwarf.ArgsMeta)
		for _, arg := range findM.Args {
			args[arg.Name] = arg
		}
		ad.pool.Set(findM.MethodName, &def.MethodInfo{
			Method:     en,
			MethodName: findM.MethodName,
			Param:      args,
		})
		initFnCache.SetDone()
		ad.hasInitPacked = true
		log.Infof("[%s] %s(%s) mapping url = %s", en.Method, ad.trimPrefix(findM.MethodName), printArgs(findM.Args), en.Url)
	})

	log.Infof("init use %s", time.Since(start))
}

func printArgs(args []dwarf.ArgsMeta) string {
	var s strings.Builder
	l := len(args) - 1
	for i, arg := range args {
		s.WriteString(arg.Name)
		if i != l {
			s.WriteString(",")
		}
	}
	return s.String()
}
