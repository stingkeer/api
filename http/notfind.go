package http

import (
	"math"
	"net/http"

	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/intercept"
)

var (
	_ intercept.HttpIntercept = (*NotFind)(nil)
)

type NotFind struct {
	serialize def.Serialize
}

func NewNotFind(serialize def.Serialize) *NotFind {
	return &NotFind{serialize: serialize}
}

// Http implements intercept.HttpIntercept.
func (n *NotFind) Http(rw http.ResponseWriter, req *http.Request, ctx *intercept.HttpContext) bool {
	// path matched a route but the method did not: 405 (+Allow) instead of a
	// silently empty 200
	if _, load := ctx.LoadAndDelete("MATCH_METHOD"); load {
		allow, _ := ctx.LoadAndDelete("MATCH_METHOD_ALLOW")
		allowStr, _ := allow.(string)
		n.methodNotAllowed(rw, req, allowStr)
		return true
	}
	if _, load := ctx.LoadAndDelete("MATCH"); load {
		n.notFindPath(rw, req, "Not find Path")
		return true
	}
	return true
}

// Order implements intercept.HttpIntercept.
func (*NotFind) Order() def.HandlerOrder {
	return math.MaxUint
}

func (api *NotFind) methodNotAllowed(rw http.ResponseWriter, req *http.Request, allow string) {
	con := api.serialize.Encode(map[string]string{
		"path":   req.URL.String(),
		"msg":    http.StatusText(http.StatusMethodNotAllowed),
		"method": req.Method,
	})
	header := rw.Header()
	if allow != "" {
		header.Set("Allow", allow)
	}
	header.Set("Content-Type", con.ContentType)
	rw.WriteHeader(http.StatusMethodNotAllowed)
	rw.Write(con.Bytes)
}

func (api *NotFind) notFindPath(rw http.ResponseWriter, req *http.Request, msg string) {
	con := api.serialize.Encode(map[string]string{
		"path": req.URL.String(),
		"msg":  "Not find Path",
	})
	header := rw.Header()
	header.Set("Content-Type", con.ContentType)
	rw.WriteHeader(http.StatusNotFound)
	rw.Write(con.Bytes)
}
