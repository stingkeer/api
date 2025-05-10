package http

import (
	"math"
	"net/http"

	"go.aew.app/api/def"
	"go.aew.app/api/intercept"
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

func (api *NotFind) notFindPath(rw http.ResponseWriter, req *http.Request, msg string) {
	con := api.serialize.Encode(map[string]string{
		"path": req.URL.String(),
		"msg":  "Not find Path",
	})
	header := rw.Header()
	header.Add("Content-Type", con.ContentType)
	rw.WriteHeader(http.StatusNotFound)
	rw.Write(con.Bytes)
}
