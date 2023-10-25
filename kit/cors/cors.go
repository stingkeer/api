package cors

import (
	"net/http"

	"gitee.com/fast_api/api"
	"gitee.com/fast_api/api/def"
	"github.com/rs/cors"
)

func init() {
	api.AddHttpHandle(NewCorsHandle(cors.Default()))
}

type Handle struct {
	cors *cors.Cors
}

func NewCorsHandle(cors *cors.Cors) *Handle {
	return &Handle{cors: cors}
}

func (c *Handle) Http(rw http.ResponseWriter, req *http.Request) bool {
	c.cors.HandlerFunc(rw, req)
	return false
}

func (c *Handle) Order() def.HandlerOrder {
	return 10
}
