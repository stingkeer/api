package rest

import (
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"os"
	"testing"

	"gitee.com/fast_api/api"
	"gitee.com/fast_api/api/cache"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/test/r"
)

func TestGetHeader(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(a def.Header) any {
			return a.Values("token")
		}, "/header")
	}).Request().AddHeader("token", "asdfed").Do(func(resp *r.Response) {
		if resp.BodyString() != "[\"asdfed\"]" {
			t.Errorf("except asdfed but %s", resp.BodyString())
		}
	})

}

func TestSetHeader(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(a def.Header) {
			a.Add("token", "nnnnn")
		}, "/no")
	}).DoRequestNobody(func(resp *r.Response) {
		if resp.Header("token") != "nnnnn" {
			t.Error("not set header")
		}
	})
}

func TestBigInt(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(a big.Int) any {
			return a.String()
		}, "/bigint")
	}).Request().AddParam("a", "10000").Do(func(resp *r.Response) {
		if resp.BodyString() != "10000" {
			t.Error("big.Int error")
		}
	})
}

func TestCacheKey(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(a int, hello def.String[cache.Key]) any {
			fmt.Println(a, hello, hello.String())
			return nil
		}, "/m")
	}).Request().AddParam("hello", "my").Do(func(resp *r.Response) {
		resp.Dump()
	})
}

func TestDownFile(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			f, e := os.Open("d:/download/QmfWv8FfpKiCWsueKfXDLrgyqXZsEuGFJFBL7TfjNmxkAw")
			fmt.Println(e)
			return api.NewStream(f).SetRateLimit(500000)
		}, "/download")
	})
}

func MulFile(read multipart.Reader) string {
	par, _ := read.NextPart()
	fmt.Println(par.FileName(), par.FormName())
	b, _ := io.ReadAll(par)
	fmt.Println(string(b))
	return string(b)
}

func TestUpload(t *testing.T) {
	api.POST(MulFile, "/update")
	api.StartService(nil)
}

func TestResp404(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			return api.NewResp(nil).SetCode(404)
		}, "/resp")
	}).DoRequestNobody(func(resp *r.Response) {
		if resp.Code() != 404 {
			t.Errorf("except 404 but %d", resp.Code())
		}
	})
}

func TestResp(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			return api.NewResp(map[string]any{
				"status": true,
			})
		}, "/resp")
	}).DoRequestNobody(func(resp *r.Response) {
		fmt.Println(resp.BodyString())
		if resp.BodyString() != "{\"status\":true}" {
			t.Error("Resp error")
		}
	})
}

func TestPostBody(t *testing.T) {
	x := "{\"name\":\"w\",\"pass\":\"12345\"}"
	r.Test(t, func() def.Option {
		return api.POST(func(a struct {
			Name string `json:"name,omitempty"`
			Pass string `json:"pass,omitempty"`
		}) any {
			return a
		}, "/login")
	}).Request().SetBody([]byte(x)).Do(func(resp *r.Response) {
		if resp.BodyString() != x {
			t.Error("TestPostBody error")
		}
	})
}
