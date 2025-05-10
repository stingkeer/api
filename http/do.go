package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"sort"
	"sync"

	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/intercept"
	"go.aew.app/api.v1/log"
)

type Handles []intercept.HttpIntercept

var (
	httpHandleZero    Handles
	httpHandles       Handles
	httpHandlesGT1000 Handles

	pool = sync.Pool{
		New: func() any {
			return intercept.NewHttpContext()
		},
	}
)

func (h Handles) Sort() {
	sort.Slice(h, func(i, j int) bool {
		return h[i].Order() < h[j].Order()
	})
}

func DoHttp(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			WriteError(handleError(err), rw)
			log.Error(err)
			debug.PrintStack()
		}
	}()

	mP := pool.Get().(*intercept.HttpContext)
	defer func() {
		mP.Clear()
		pool.Put(mP)
	}()

	//execute order == 0 handler
	for _, handle := range httpHandleZero {
		if handle != nil {
			if handle.Http(rw, req, mP) {
				break
			}
		}
	}

	//execute system and user handler
	for _, handle := range httpHandles {
		if handle != nil {
			if handle.Http(rw, req, mP) {
				break
			}
		}
	}

	if mP.IsSkipResponse() {
		return
	}

	//Response area
	//execute order >= 1000 handler
	for _, handle := range httpHandlesGT1000 {
		if handle != nil {
			if handle.Http(rw, req, mP) {
				break
			}
		}
	}
}

func addHttpHandle(f intercept.HttpIntercept) {
	if f.Order() == 0 {
		httpHandleZero = append(httpHandleZero, f)
		return
	}
	if f.Order() >= 1000 {
		httpHandlesGT1000 = append(httpHandlesGT1000, f)
		httpHandlesGT1000.Sort()
		return
	}
	httpHandles = append(httpHandles, f)
	httpHandles.Sort()
}

func AddHttpHandle(f intercept.HttpIntercept) {
	if f.Order() == 0 {
		addHttpHandle(f)
		return
	}
	if f.Order() > 100 && f.Order() < 1000 {
		panic(fmt.Errorf("HttpIntercept order %d Must be greater than or equal to 100", f.Order()))
	}
	addHttpHandle(f)
}

func WriteError(err any, rw http.ResponseWriter) {
	bytes, e := json.Marshal(err)
	if e != nil {
		panic(e)
	}
	rw.Header().Add(def.Content_Type, def.Content_JSON)
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write(bytes)
}
