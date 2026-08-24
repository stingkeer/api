package rest

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"go.aew.app/api.v1"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/test/r"
)

func TestGzip(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			return map[string]string{"status": "OK"}
		}, "/gzip")
	}).Request().AddHeader("Accept-Encoding", "gzip").Do(func(resp *r.Response) {
		resp.AssertHeader("Content-Encoding", "gzip")
		body := gunzipBody(t, resp)
		var m map[string]string
		if err := json.Unmarshal(body, &m); err != nil {
			t.Fatalf("unmarshal %s: %v", body, err)
		}
		if m["status"] != "OK" {
			t.Errorf("expect %q but %q", "OK", m["status"])
		}
	})
}

func TestHtmlGzip(t *testing.T) {
	type H struct {
		Hello string
	}
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			return api.Html(`<h>{{.Hello}}</h>`, H{Hello: "my"})
		}, "/gzip-html")
	}).Request().AddHeader("Accept-Encoding", "gzip").Do(func(resp *r.Response) {
		resp.AssertHeader("Content-Encoding", "gzip")
		body := gunzipBody(t, resp)
		if !strings.Contains(string(body), "<h>my</h>") {
			t.Errorf("expect <h>my</h> in %q", body)
		}
	})
}

// gunzipBody decompresses a gzip response body.
func gunzipBody(t *testing.T, resp *r.Response) []byte {
	t.Helper()
	zr, err := gzip.NewReader(bytes.NewReader(resp.Body()))
	if err != nil {
		t.Fatalf("gzip reader: %v", err)
	}
	defer zr.Close()
	bys, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("read gunzipped body: %v", err)
	}
	return bys
}
