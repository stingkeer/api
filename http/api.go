package http

import (
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/intercept"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type ApiInter struct {
	match     def.Match
	caller    def.Caller
	serialize def.Serialize
}

func NewApiIntercept(match def.Match, caller def.Caller, serialize def.Serialize) intercept.HttpIntercept {
	return &ApiInter{
		match:     match,
		caller:    caller,
		serialize: serialize,
	}
}

func (api *ApiInter) Http(rw http.ResponseWriter, req *http.Request) bool {
	logrus.Tracef("incoming req Method [%s] , Url [%s]", req.Method, req.URL.String())
	entry := api.match.Match(req.URL)
	req.Header.Del(def.HEAD_CONST)
	if nil == entry {
		logrus.Tracef("not match %s", req.URL)
		return false
	}

	if req.Method != entry.Method {
		logrus.Warnf("not support Method %s", req.Method)
		return false
	}
	if entry.Fn != nil {

		iRet := api.caller.Call(entry, req)

		if !doWithRet(iRet, rw, req) {
			if iRet == nil {
				WriteResponse(rw, req, nil)
				return false
			}

			//default return json
			h := api.serialize.Encode(iRet)
			if h == nil {
				WriteResponse(rw, req, nil)
				return false
			}
			WriteResponse(rw, req, h)
			return false
		}

	}
	return false
}

func doWithRet(value interface{}, rw http.ResponseWriter, req *http.Request) bool {
	typ := reflect.TypeOf(value)
	if _, b := retAdapters[typ]; b {
		WriteRetResponse(rw, req, value.(def.RetAdapter))
		return true
	} else {
		return false
	}
}

func WriteRetResponse(rw http.ResponseWriter, req *http.Request, adapter def.RetAdapter) {
	appendSysHeader(rw, req)
	header := rw.Header()
	header.Add("Content-Type", adapter.ContentType())

	//if struct impl def.AppendHeader ,can def header
	if v, b := adapter.(def.AppendHeader); b {
		mHeader := v.Append(&readHead{
			req: req,
		})
		for k, v := range mHeader {
			header.Add(k, v)
		}
	}

	if v, b := adapter.(def.HttpStatus); b {
		rw.WriteHeader(v.Code())
	}

	_, err := io.Copy(rw, adapter.Return())
	if err != nil {
		logrus.Error(err)
	}

	if v, b := adapter.(io.Closer); b {
		v.Close()
	}

}

func appendSysHeader(rw http.ResponseWriter, req *http.Request) {
	header := rw.Header()
	apHeader := req.Header.Get(def.HEAD_CONST)
	if apHeader != "" {
		kvs := strings.Split(apHeader, ",")
		for _, v := range kvs {
			l := strings.Split(v, "=")
			if len(l) != 2 {
				panic("you set head error")
			}
			header.Add(l[0], l[1])
		}
	}
}

func WriteResponse(rw http.ResponseWriter, req *http.Request, content *def.Content) {
	appendSysHeader(rw, req)
	header := rw.Header()
	if content != nil {
		header.Add("Content-Type", content.ContentType)
		rw.WriteHeader(http.StatusOK)
		if content.Bytes != nil {
			rw.Write(content.Bytes)
		}
	}

}

func (api *ApiInter) Order() int {
	return 100
}
