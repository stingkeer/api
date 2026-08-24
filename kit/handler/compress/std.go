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
	// Range semantics apply to the representation bytes; re-encoding the
	// ranged body (different byte offsets, different length) corrupts
	// Content-Range/Content-Length. Serve ranged responses uncompressed.
	if req.Header.Get("Range") != "" {
		return false
	}
	if cmp := g.checkSupport(req.Header.Get(def.Accept_Encoding)); cmp != nil {
		if c, b := ctx.LoadAndDelete("CALLDATA_RetAdapter"); b {
			src := c.(def.RetAdapter)

			// Collect status/headers/content-type from the source adapter BEFORE
			// starting the copy goroutine: Append/ContentType mutate shared state
			// (seek the underlying reader, set range bounds) and racing them
			// against the copy goroutine is a data race on Stream fields and the
			// underlying reader.
			code := http.StatusOK
			if v, is := c.(def.HttpStatus); is {
				code = v.Code()
			}
			mHeader := make(map[string]string)
			if v, is := c.(def.AppendHeader); is {
				for k, v1 := range v.Append(&readHead{req: req}) {
					// compressed byte count differs from source length; a stale
					// Content-Length truncates the response mid-decompression
					if k == "Content-Length" {
						continue
					}
					mHeader[k] = v1
				}
			}
			ct := src.ContentType()

			r, w := io.Pipe()
			target := cmp.New(w)
			go func() {
				defer w.Close()
				if _, err := io.Copy(target, src.Return()); err != nil {
					return
				}
				if err := target.Close(); err != nil {
					return
				}
			}()

			resp := rettypes.NewStream(r)
			resp.SetCode(code)
			for k, v1 := range mHeader {
				resp.AddHeader(k, v1)
			}
			resp.SetContentType(ct)
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
