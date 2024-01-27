package http

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"reflect"

	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/intercept"
	"gitee.com/fast_api/api/log"
)

var (
	_ intercept.HttpIntercept = (*ApiInter)(nil)
	_ intercept.HttpIntercept = (*ApiRespose)(nil)
)

type ApiInter struct {
	match     def.Match
	caller    def.Caller
	serialize def.Serialize
	pool      *def.MethodsPools
}

type ApiRespose struct {
}

func NewApiRespose() *ApiRespose {
	return &ApiRespose{}
}

// Http implements intercept.HttpIntercept.
func (resp *ApiRespose) Http(rw http.ResponseWriter, req *http.Request, ctx *intercept.HttpContext) bool {
	if _, load := ctx.LoadAndDelete("CALLDATA_NIL"); load {
		WriteResponse(rw, req, nil)
		return true
	}
	// RetAdapter handling
	if v, load := ctx.LoadAndDelete("CALLDATA_RetAdapter"); load {
		WriteRetResponse(rw, req, v.(def.RetAdapter))
		return true
	}

	if v, load := ctx.LoadAndDelete("CALLDATA"); load {
		WriteResponse(rw, req, v.(*def.Content))
		return true
	}
	fmt.Printf("no match %p %s\n", ctx, req.RequestURI)
	return false
}

// Order implements intercept.HttpIntercept.
func (*ApiRespose) Order() def.HandlerOrder {
	return math.MaxUint - 100
}

func NewApiIntercept(match def.Match, caller def.Caller, serialize def.Serialize, pool *def.MethodsPools) intercept.HttpIntercept {
	return &ApiInter{
		match:     match,
		caller:    caller,
		serialize: serialize,
		pool:      pool,
	}
}

func (api *ApiInter) Http(rw http.ResponseWriter, req *http.Request, ctx *intercept.HttpContext) bool {
	log.Tracef("incoming req HttpMethod [%s] , Url [%s]", req.Method, req.URL.String())
	entry := api.match.Match(req.URL)
	req.Header.Del(def.HEAD_CONST)
	if nil == entry {
		log.Tracef("not match %s", req.URL)
		ctx.Store("Match", 0)
		return false
	}
	if req.Method != entry.HttpMethod {
		log.Warnf("not support HttpMethod %s", req.Method)
		ctx.Store("Match_Method", req.Method)
		return false
	}
	if entry.Fn != nil {
		iRet := api.caller.Call(entry, def.WithRequest(rw, req))
		//Returns null handling
		if iRet == def.Empty("") {
			return false
		}
		if iRet == nil {
			ctx.Store("CALLDATA_NIL", 1)
			return false
		}
		// RetAdapter handling
		typ := reflect.TypeOf(iRet)
		if _, b := retAdapters[typ]; b {
			ctx.Store("CALLDATA_RetAdapter", iRet)
			return false
		}
		// Serialize the value
		h := api.serialize.Encode(iRet)
		if h != nil {
			ctx.Store("CALLDATA", h)
			return false
		}
		return false

	}
	return false
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
