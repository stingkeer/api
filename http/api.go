package http

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/intercept"
	"gitee.com/fast_api/api/log"
)

type ApiInter struct {
	match     def.Match
	caller    def.Caller
	serialize def.Serialize
	pool      *def.MethodsPools
}

func NewApiIntercept(match def.Match, caller def.Caller, serialize def.Serialize, pool *def.MethodsPools) intercept.HttpIntercept {
	return &ApiInter{
		match:     match,
		caller:    caller,
		serialize: serialize,
		pool:      pool,
	}
}

func (api *ApiInter) Http(rw http.ResponseWriter, req *http.Request) bool {
	log.Tracef("incoming req HttpMethod [%s] , Url [%s]", req.Method, req.URL.String())
	entry := api.match.Match(req.URL)
	req.Header.Del(def.HEAD_CONST)
	if nil == entry {
		log.Tracef("not match %s", req.URL)
		return false
	}
	if req.Method != entry.HttpMethod {
		log.Warnf("not support HttpMethod %s", req.Method)
		return false
	}
	if entry.Fn != nil {
		iRet := api.caller.Call(entry, req)
		if iRet == nil {
			WriteResponse(rw, req, nil)
			return true
		}
		if doWithRet(iRet, rw, req) {
			return true
		}
		h := api.serialize.Encode(iRet)
		if h != nil {
			WriteResponse(rw, req, h)
			return true
		}
		return false

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
		for k, v1 := range mHeader {
			header.Add(k, v1)
		}
	}

	if v, b := adapter.(def.HttpStatus); b {
		rw.WriteHeader(v.Code())
	}

	_, err := io.Copy(rw, adapter.Return())
	if err != nil {
		log.Error(err)
	}

	if v, b := adapter.(io.Closer); b {
		v.Close()
	}

}

func appendSysHeader(rw http.ResponseWriter, req *http.Request) {
	header := rw.Header()
	apHeader := req.Header.Get(def.HEAD_CONST)
	if apHeader != "" {
		m := make(map[string]string)
		if err := json.Unmarshal([]byte(apHeader), &m); err == nil {
			for k, v := range m {
				header.Add(k, v)
			}
		}
	}
}

func WriteResponse(rw http.ResponseWriter, req *http.Request, content *def.Content) {
	appendSysHeader(rw, req)
	header := rw.Header()
	if content != nil {
		header.Add("Content-Type", content.ContentType)
		rw.WriteHeader(http.StatusOK)
		rw.Write(content.Bytes)
	}

}

func (api *ApiInter) Order() def.HandlerOrder {
	return def.Handler_API
}
