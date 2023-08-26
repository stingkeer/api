package cors

import (
	"gitee.com/fast_api/api"
	"github.com/rs/cors"
	"net/http"
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

func (c *Handle) Order() int {
	return 10
}
