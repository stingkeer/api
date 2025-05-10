package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"os"
	"testing"

	"go.aew.app/api.v1"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/test/r"
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

func TestInt64(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(a int64, b int64) any {
			return a
		}, "/int64")
	}).Request().AddParam("a", "1").Do(func(resp *r.Response) {
		fmt.Println(resp.BodyString())
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

func TestBodyString(t *testing.T) {
	var ss string
	fmt.Println(json.Unmarshal([]byte("dasdf"), &ss))
	//because many options
	//Post /body?a=xxx
	//Post /body raw data
	r.Test(t, func() def.Option {
		return api.POST(func(body string) any {
			return body
		}, "/body")
	}).Request().SetBody([]byte("hello")).Do(func(resp *r.Response) {
		fmt.Println(resp.BodyString())
	})
}

func TestBigIntPtr(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(a *big.Int) any {
			return a.String()
		}, "/bigint")
	}).Request().AddParam("a", "10000").Do(func(resp *r.Response) {
		if resp.BodyString() != "10000" {
			t.Error("big.Int error")
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

func TestRequirParam(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(name def.StringReq) {

		}, "/ased")
	}).Request().Do(func(resp *r.Response) {
		var e def.Error
		if er := json.Unmarshal(resp.Body(), &e); er != nil {
			t.Error(er)
		}
		if e.Code != 3808709076 {
			t.Error("code err")
		}
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
		resp.AssetBody("{\"status\":true}")
	})
}

func TestPostBody(t *testing.T) {
	x := "{\"name\":\"w\",\"pass\":\"12345\"}"
	r.Test(t, func() def.Option {
		return api.POST(func(body struct {
			Name string `json:"name,omitempty"`
			Pass string `json:"pass,omitempty"`
		}) any {
			return body
		}, "/login")
	}).Request().SetBody([]byte(x)).Do(func(resp *r.Response) {
		resp.AssetBody(x)
	})
}

func TestBodyBytes(t *testing.T) {
	var ss string
	fmt.Println(json.Unmarshal([]byte("dasdf"), &ss))
	//because many options
	//Post /body?a=xxx
	//Post /body raw data
	r.Test(t, func() def.Option {
		return api.POST(func(body []byte) any {
			return body
		}, "/body")
	}).Request().SetBody([]byte("hello")).Do(func(resp *r.Response) {
		fmt.Println(resp.BodyString())
	})
}

func TestBodySlice(t *testing.T) {
	str := `["aaa","bbbb","ccc"]`
	//because many options
	//Post /body?a=xxx
	//Post /body raw data
	r.Test(t, func() def.Option {
		return api.POST(func(body []string) any {
			return body
		}, "/body")
	}).Request().SetBody([]byte(str)).Do(func(resp *r.Response) {
		if str != resp.BodyString() {
			t.Error("TestBodySlice fail")
		}
	})
}

func TestBodyAndParam(t *testing.T) {
	x := "{\"name\":\"w\",\"pass\":\"12345\"}"
	r.Test(t, func() def.Option {
		return api.POST(func(body struct {
			Name string `json:"name,omitempty"`
			Pass string `json:"pass,omitempty"`
		}, ok def.StringReq) any {
			fmt.Println("ok = ", ok)
			return body
		}, "/bodyParam")
	}).Request().SetBody([]byte(x)).AddParam("ok", "hello").Do(func(resp *r.Response) {
		resp.AssetBody(x)
	})
}

func TestOnlyParam1(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.POST(func(aaa struct {
			Name string `json:"name,omitempty"`
			Pass string `json:"pass,omitempty"`
		}) any {
			return aaa
		}, "/bodyParam")
	}).Request().AddParam("name", "hello").AddParam("pass", "2342342").Do(func(resp *r.Response) {
		fmt.Println(resp.BodyString())
	})
}

func TestOnlyParam2(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.POST(func(aaa struct {
			Name string `json:"name1,omitempty"`
			Pass string `json:"pass,omitempty"`
		}, bbb struct {
			Name string `json:"name,omitempty"`
			Pass string `json:"pass,omitempty"`
		}) any {
			return map[string]any{"bbbb": bbb, "aaa": aaa}
		}, "/onlyParam")
	}).Request().AddParam("name", "hello").AddParam("pass", "2342342").AddParam("name1", "hello").Do(func(resp *r.Response) {
		fmt.Println(resp.BodyString())
	})
}

func TestOnlyParam3(t *testing.T) {
	type Page struct {
		OffSet int `json:"offSet"`
		Page   int `json:"page"`
	}
	r.Test(t, func() def.Option {
		return api.POST(func(aaa struct {
			Page
			Name string `json:"name1,omitempty"`
			Pass string `json:"pass,omitempty"`
		}) any {
			return map[string]any{"aaa": aaa}
		}, "/onlyParam")
	}).Request().
		AddParam("name", "hello").
		AddParam("pass", "2342342").
		AddParam("pass", "2342342").
		AddParam("name1", "hello").
		AddParam("offSet", "1").
		AddParam("page", "10").
		Do(func(resp *r.Response) {
			fmt.Println(resp.BodyString())
		})
}
