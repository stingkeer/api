package http

import (
	"encoding/json"
	"gitee.com/fast_api/api/intercept"
	"gitee.com/fast_api/api/public"
	"github.com/sirupsen/logrus"
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
			debug.PrintStack()
			logrus.Error(err)
			if v, b := err.(string); b {
				WriteError(public.NewError(v), rw)
			}
			if v, b := err.(error); b {
				WriteError(public.NewError(v.Error()), rw)
			}
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

func WriteError(err public.Error, rw http.ResponseWriter) {
	rw.Header().Add("Content-Type", public.Json)
	rw.WriteHeader(500)
	bytes, _ := json.Marshal(err)
	rw.Write(bytes)
}
