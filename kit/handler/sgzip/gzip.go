package sgzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"go.aew.app/api/call/rettypes"
	"go.aew.app/api/def"
	"go.aew.app/api/intercept"
)

var _ intercept.HttpIntercept = (*GZip)(nil)

type GZip struct{}

type readHead struct {
	req *http.Request
}

func (r *readHead) Get(key string) string {
	return r.req.Header.Get(key)
}

func (r *readHead) Values(key string) []string {
	return r.req.Header.Values(key)
}

// Http implements intercept.HttpIntercept.
func (g *GZip) Http(rw http.ResponseWriter, req *http.Request, ctx *intercept.HttpContext) bool {
	if strings.Contains(req.Header.Get(def.Accept_Encoding), "gzip") {
		if c, b := ctx.LoadAndDelete("CALLDATA_RetAdapter"); b {
			g := c.(def.RetAdapter)
			r, w := io.Pipe()
			go func() {
				gz := gzip.NewWriter(w)
				_, err := io.Copy(gz, g.Return())
				if err != nil {
					w.Close()
					return
				}
				if err := gz.Close(); err != nil {
					w.Close()
					return
				}
				w.Close()
			}()
			resp := rettypes.NewStream(r)

			if status, is := c.(def.HttpStatus); is {
				resp.SetCode(status.Code())
			}
			//if struct impl def.AppendHeader ,can def header
			if v, b := c.(def.AppendHeader); b {
				//Call the return handler and assign append
				mHeader := v.Append(&readHead{
					req: req,
				})
				for k, v1 := range mHeader {
					resp.AddHeader(k, v1)
				}
			}
			resp.SetContentType(g.ContentType())
			resp.AddHeader(def.Content_Encoding, "gzip")
			ctx.Store("CALLDATA_RetAdapter", resp)
		}

		if c, b := ctx.LoadAndDelete("CALLDATA"); b {
			g := c.(*def.Content)
			buf := new(bytes.Buffer)
			gw := gzip.NewWriter(buf)
			_, err := gw.Write(g.Bytes)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return true
			}
			if err := gw.Close(); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return true
			}
			resp := rettypes.NewStream(buf)
			resp.SetContentType(g.ContentType)
			resp.AddHeader(def.Content_Encoding, "gzip")
			ctx.Store("CALLDATA_RetAdapter", resp)
		}
	}
	return false
}

// Order implements intercept.HttpIntercept.
func (g *GZip) Order() def.HandlerOrder {
	return def.Handler_GZIP
}
