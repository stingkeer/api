package http

import (
	"encoding/json"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/intercept"
	"gitee.com/fast_api/api/log"
	"net/http"
	"runtime/debug"
	"sort"
	"sync"
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

func AddHttpHandle(f intercept.HttpIntercept) {
	httpHandles = append(httpHandles, f)
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
