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
	r "gitee.com/fast_api/api/test/R"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}

func TestGetHeader(t *testing.T) {
	r.Test(func() {
		api.GET(func(a def.Header) any {
			return a.Values("Accept-Encoding")
		}, "/h")
	})
}

func TestSetHeader(t *testing.T) {
	r.Test(func() {
		api.GET(func(a def.Header) {
			a.Add("szb", "nnnnn")
		}, "/no")
	})
}

func TestBigInt(t *testing.T) {
	r.Test(func() {
		api.GET(func(a big.Int) any {
			return a.String()
		}, "/int")
	})
}

func TestCacheKey(t *testing.T) {
	r.Test(func() {
		api.GET(func(a int, hello def.String[cache.Key]) any {
			fmt.Println(a, hello, hello.String())
			return nil
		}, "/m")
	})
}

func TestDownFile(t *testing.T) {
	r.Test(func() {
		api.GET(func() any {
			f, e := os.Open("d:/download/QmfWv8FfpKiCWsueKfXDLrgyqXZsEuGFJFBL7TfjNmxkAw")
			fmt.Println(e)
			stream := api.NewStream(f)
			stream.SetRateLimit(500000)
			return stream
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
	r.Test(func() {
		api.GET(func() any {
			return api.NewResp(nil).SetCode(404)
		}, "/resp")
	})
}

func TestResp(t *testing.T) {
	r.Test(func() {
		api.GET(func() any {
			return api.NewResp(map[string]any{
				"status": true,
			})
		}, "/resp")
	})
}
