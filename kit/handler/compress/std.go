package compress

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"go.aew.app/api.v1/call/rettypes"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/intercept"
)

var CompressRegister = map[string]Compress{
	"gzip":    &gZip{},
	"deflate": &flateStd{},
}

var _ intercept.HttpIntercept = (*CompressStd)(nil)

type CompressStd struct{}

type readHead struct {
	req *http.Request
}

func (r *readHead) Get(key string) string {
	return r.req.Header.Get(key)
}

func (r *readHead) Values(key string) []string {
	return r.req.Header.Values(key)
}

// checkSupport
// gzip, deflate, br, zstd
func (g *CompressStd) checkSupport(h string) Compress {
	for k, v := range CompressRegister {
		if strings.Contains(h, k) {
			return v
		}
	}
	return nil
}

// Http implements intercept.HttpIntercept.
func (g *CompressStd) Http(rw http.ResponseWriter, req *http.Request, ctx *intercept.HttpContext) bool {
	if cmp := g.checkSupport(req.Header.Get(def.Accept_Encoding)); cmp != nil {
		if c, b := ctx.LoadAndDelete("CALLDATA_RetAdapter"); b {
			g := c.(def.RetAdapter)
			r, w := io.Pipe()
			target := cmp.New(w)
			go func() {
				_, err := io.Copy(target, g.Return())
				if err != nil {
					w.Close()
					return
				}
				if err := target.Close(); err != nil {
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
			resp.AddHeader(def.Content_Encoding, cmp.ContentEncoding())
			ctx.Store("CALLDATA_RetAdapter", resp)
		}

		if c, b := ctx.LoadAndDelete("CALLDATA"); b {
			g := c.(*def.Content)
			buf := new(bytes.Buffer)
			gw := cmp.New(buf)
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
			resp.AddHeader(def.Content_Encoding, cmp.ContentEncoding())
			ctx.Store("CALLDATA_RetAdapter", resp)
		}
	}
	return false
}

// Order implements intercept.HttpIntercept.
func (g *CompressStd) Order() def.HandlerOrder {
	return def.Handler_HTTP_COMPRESS
}
