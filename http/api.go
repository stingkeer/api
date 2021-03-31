package http

import (
	"gitee.com/fast_api/api/public"
	"gitee.com/fast_api/api/server"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ApiInter struct {
	match     public.Match
	caller    public.Caller
	serialize public.Serialize
}

func (api *ApiInter) Http(rw http.ResponseWriter, req *http.Request) bool {

	logrus.Tracef("incoming req Method [%s] , Url [%s]", req.Method, req.URL.String())
	entry := api.match.Match(req.URL)
	if nil == entry {
		logrus.Tracef("not match %s", req.URL)
		return false
	}

	if req.Method != entry.Method {
		logrus.Warnf("not support Method %s", req.Method)
		return false
	}

	if entry.Fn != nil {
		inf := api.caller.Call(entry, req)
		if inf == nil {
			return false
		}
		h := api.serialize.Encode(inf)
		if h == nil {
			return false
		}
		rw.Header().Add("Content-Type", h.ContentType)
		rw.WriteHeader(h.Code)
		rw.Write(h.Bytes)
		return false
	}
	return false
}

func (api *ApiInter) Order() int {
	return 100
}

func init() {
	server.Invoke(func(match public.Match, resultConvert public.Serialize, caller public.Caller) {
		AddHttpHandle(&ApiInter{
			match,
			caller,
			resultConvert,
		})
	})
}
