package rest

import (
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"gitee.com/fast_api/api"
	"gitee.com/fast_api/api/cache"
	"gitee.com/fast_api/api/def"
	"github.com/sirupsen/logrus"
)

func TestApiHttp(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	api.GET(func(a http.Request) (any, any) {
		return a.Host, cache.NewCacheImpl(time.Second)
	}, "/hello")

	api.GET(func(a big.Int) any {
		return a.String()
	}, "/int")

	api.GET(func(a def.Header) any {
		return a.Values("Accept-Encoding")
	}, "/h")

	api.GET(func(a def.Header) {
		a.Add("szb", "nnnnn")
	}, "/no")

	api.GET(func(a int, hello def.String[cache.Key]) any {
		fmt.Println(a, hello, hello.String())
		return nil
	}, "/m")

	api.GET(func(req *http.Request, resp http.Response) {

	}, "/http")

	api.GET(func() any {
		f, e := os.Open("d:/download/QmfWv8FfpKiCWsueKfXDLrgyqXZsEuGFJFBL7TfjNmxkAw")
		fmt.Println(e)
		stream := api.NewStream(f)
		stream.SetRateLimit(500000)
		return stream
	}, "/download")

	api.GET(func(header def.Header, reader multipart.Reader) {
		fmt.Println(header, reader)
	}, "/file")

	api.StartService(nil)
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
