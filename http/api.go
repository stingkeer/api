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
		iRet := api.caller.Call(entry, def.WithRequest(rw, req))
		//Returns null handling
		if iRet == nil {
			// TODO Fix me
			// WriteResponse(rw, req, nil)
			return true
		}
		// RetAdapter handling
		if doWithRet(iRet, rw, req) {
			return true
		}
		// Serialize the value
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
	//set ContentType
	header.Add("Content-Type", adapter.ContentType())

	//if struct impl def.AppendHeader ,can def header
	if v, b := adapter.(def.AppendHeader); b {
		//Call the return handler and assign append
		mHeader := v.Append(&readHead{
			req: req,
		})
		for k, v1 := range mHeader {
			header.Add(k, v1)
		}
	}

	//set http status
	if v, b := adapter.(def.HttpStatus); b {
		rw.WriteHeader(v.Code())
	}

	//return http io
	_, err := io.Copy(rw, adapter.Return())
	if err != nil {
		log.Error(err)
	}

	//close
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
