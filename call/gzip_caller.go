package call

import (
	"bytes"
	"compress/gzip"
	"strings"

	"gitee.com/fast_api/api/call/rettypes"
	"gitee.com/fast_api/api/def"
)

var _ def.Caller = (*GzipCaller)(nil)

type GzipCaller struct {
	WsCaller
}

func NewGzipCaller(serialize def.Serialize, pool *def.MethodsPools) *GzipCaller {
	return &GzipCaller{
		WsCaller: *NewWsCaller(serialize, pool),
	}
}

// Call implements def.Caller.
func (g *GzipCaller) Call(f *def.Entry, req *def.Request) interface{} {
	call := g.WsCaller.Call(f, req)
	if call != nil && strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
		//"Content-Encoding", "gzip"
		buf := new(bytes.Buffer)
		gw := gzip.NewWriter(buf)
		ctx := g.serialize.Encode(call)
		if _, err := gw.Write(ctx.Bytes); err != nil {
			panic(err)
		}
		if err := gw.Close(); err != nil {
			panic(err)
		}
		return rettypes.NewResp(call).SetReader(buf).SetHeader(map[string]string{def.Content_Encoding: "gzip"})
	}
	return call
}
