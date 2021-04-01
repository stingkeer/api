package http

import (
	"gitee.com/fast_api/api/public"
	"gitee.com/fast_api/api/server"
	"github.com/sirupsen/logrus"
	"net/http"
	"reflect"
	"strings"
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
			WriteResponse(rw, req, nil)
			return false
		}
		h := api.serialize.Encode(inf)
		if h == nil {
			WriteResponse(rw, req, nil)
			return false
		}
		WriteResponse(rw, req, h)
		return false
	}
	return false
}

func WriteResponse(rw http.ResponseWriter, req *http.Request, content *public.Content) {
	header := rw.Header()
	apHeader := req.Header.Get(public.HEAD_CONST)
	if apHeader == "" {
		return
	}
	kvs := strings.Split(apHeader, ",")
	for _, v := range kvs {
		l := strings.Split(v, "=")
		if len(l) != 2 {
			panic("you set head error")
		}
		header.Add(l[0], l[1])
	}
	if content != nil {
		header.Add("Content-Type", content.ContentType)
		rw.Write(content.Bytes)
	}
	rw.WriteHeader(http.StatusOK)

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

	errorsMap = make(map[reflect.Type]ErrorHandler)
}
