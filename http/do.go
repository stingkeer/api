package http

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"sort"
	"sync"

	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/intercept"
	"gitee.com/fast_api/api/log"
)

type Handles []intercept.HttpIntercept

var (
	g           sync.Once
	httpHandles Handles
)

func DoHttp(rw http.ResponseWriter, req *http.Request) {
	g.Do(func() {
		sort.Slice(httpHandles, func(i, j int) bool {
			return httpHandles[i].Order() < httpHandles[j].Order()
		})
	})
	defer func() {
		if err := recover(); err != nil {
			WriteError(handleError(err), rw)
			log.Error(err)
			debug.PrintStack()
		}
	}()
	for _, handle := range httpHandles {
		if handle != nil {
			if handle.Http(rw, req) {
				break
			}
		}
	}
}

func addHttpHandle(f intercept.HttpIntercept) {
	httpHandles = append(httpHandles, f)
}

func AddHttpHandle(f intercept.HttpIntercept) {
	if f.Order() <= 100 {
		log.Error("HttpIntercept Must be greater than or equal to 100")
		return
	}
	addHttpHandle(f)
}

func WriteError(err interface{}, rw http.ResponseWriter) {
	bytes, e := json.Marshal(err)
	if e != nil {
		panic(e)
	}
	rw.Header().Add("Content-Type", def.Content_JSON)
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write(bytes)
}
