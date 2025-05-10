package cors

import (
	"net/http"

	"github.com/rs/cors"
	"go.aew.app/api/def"
)

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
	return 0
}
