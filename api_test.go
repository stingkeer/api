package api

import (
	"fmt"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/server"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func hello1(kk def.StringReq) interface{} {
	return map[string]string{"name": kk.String()}
}

func MulFile(read multipart.Reader) string {
	par, _ := read.NextPart()
	fmt.Println(par.FileName(), par.FormName())
	b, _ := ioutil.ReadAll(par)
	fmt.Println(string(b))
	return string(b)
}

func TestApiHttp(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	GET(func(a http.Request) interface{} {
		return a.Host
	}, "/hello")

	GET(func(a big.Int) interface{} {
		return a.String()
	}, "/int")

	GET(func(a def.Header) interface{} {
		return a.Values("Accept-Encoding")
	}, "/h")

	GET(func(a def.Header) {
		a.Add("szb", "nnnnn")
	}, "/no")

	GET(func(a int, hello def.StringReq) {
		fmt.Println(a, hello)
	}, "/m")

	GET(func() interface{} {
		f, e := os.Open("d:/download/QmfWv8FfpKiCWsueKfXDLrgyqXZsEuGFJFBL7TfjNmxkAw")
		fmt.Println(e)
		stream := NewStream(f)
		stream.SetRateLimit(500000)
		return stream
	}, "/download")

	GET(func(header def.Header, reader multipart.Reader) {
		fmt.Println(header, reader)
	}, "/file")

	StartService(":8011")
}

func TestURL(t *testing.T) {
	GET(hello1, "/s")
	GET(hello1, "/s/<kk>")
	POST(MulFile, "/update")
	StartService(":8080")
	//reflect.ValueOf(a.show)
}

func TestPath(t *testing.T) {
	GET(hello1, "/s/<kk>")
	StartService(":8080")
}

type A struct {
}

func (A) Encode(interface{}) *def.Content {
	return nil
}
func (A) Decode([]byte, interface{}) error {
	return nil
}

func TestFile(t *testing.T) {
	server.Provide(func() def.Serialize {
		return &A{}
	})
	POST(MulFile, "/update")
	StartService(":8080")
}

func TestType(t *testing.T) {
	//bytes 64448694-86052627/86052628
	//21603934
	//13f63f85d5b618b2c03121c0a8d38758
	//13f63f85d5b618b2c03121c0a8d38758
	//2-9
}
