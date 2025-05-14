package rest

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	}).With(func() {
		gzipClient(t, func(s []byte) {
			var r map[string]string
			if json.Unmarshal(s, &r) == nil && r["status"] != "OK" {
				t.Errorf("except %s but %s", "OK", string(s))
			}
		})
	})
}

func TestHtmlGzip(t *testing.T) {
	type H struct {
		Hello string
	}
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			return api.Html(`<h>{{.Hello}}</h>`, H{Hello: "my"})
		}, "/gzip")
	}).With(func() {
		gzipClient(t, func(s []byte) {
			fmt.Println(string(s))
		})
	})
}

func gzipClient(t *testing.T, f func(s []byte)) {
	req, err := http.NewRequest("GET", "http://localhost:8080/gzip", nil)
	if err != nil {
		fmt.Println("The creation request failed:", err)
		return
	}
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("The send request failed:", err)
		return
	}
	defer resp.Body.Close()
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			t.Error("Failed to create a gzip decompressor:", err)
			return
		}
		defer reader.Close()
		bys, err := io.ReadAll(reader)
		if err != nil {
			t.Error(err)
		}
		if f != nil {
			f(bys)
		}
	} else {
		t.Errorf("Not gzip request")
	}
}
