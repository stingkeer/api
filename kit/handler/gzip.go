package handler

import (
	"compress/gzip"
	"net/http"
	"strings"

	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/intercept"
)

var _ intercept.HttpIntercept = (*GZip)(nil)

type GZip struct{}

// Http implements intercept.HttpIntercept.
func (g *GZip) Http(rw http.ResponseWriter, req *http.Request) bool {
	if strings.Contains(req.Header.Get(def.Content_Encoding), "gzip") {
		reader, err := gzip.NewReader(req.Body)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return true
		}
		req.Body = reader
	}
	return false
}

// Order implements intercept.HttpIntercept.
func (g *GZip) Order() def.HandlerOrder {
	return def.Handler_GZIP
}
