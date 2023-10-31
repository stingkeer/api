package rest

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"gitee.com/fast_api/api"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/test/r"
)

func TestGzip(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			return map[string]string{"status": "OK"}
		}, "/gzip")
	}).With(func() {
		gzipClient(t)
	})

}

func gzipClient(t *testing.T) {
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
		var r map[string]string
		if json.Unmarshal(bys, &r) == nil && r["status"] != "OK" {
			t.Errorf("except %s but %s", "OK", string(bys))
		}
	} else {
		t.Errorf("Not gzip request")
	}
}
