package http

import (
	"net/http"

	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/intercept"
)

var _ intercept.HttpIntercept = (*NotFind)(nil)

type NotFind struct {
	serialize def.Serialize
}

func NewNotFind(serialize def.Serialize) *NotFind {
	return &NotFind{serialize: serialize}
}

// Http implements intercept.HttpIntercept.
func (n *NotFind) Http(rw http.ResponseWriter, req *http.Request) bool {
	n.notFindPath(rw, req)
	return true
}

// Order implements intercept.HttpIntercept.
func (*NotFind) Order() def.HandlerOrder {
	return def.Handler_NOTFIND
}

func (api *NotFind) notFindPath(rw http.ResponseWriter, req *http.Request) {
	con := api.serialize.Encode(map[string]string{
		"path": req.URL.String(),
		"msg":  "Not find Path",
	})
	header := rw.Header()
	header.Add("Content-Type", con.ContentType)
	rw.WriteHeader(http.StatusNotFound)
	rw.Write(con.Bytes)
}
