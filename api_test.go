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
	"reflect"
	"testing"
)

func hello1(kk string) interface{} {
	return map[string]string{"name": kk}
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

	GET(func(a int, b def.StringReq) {
		fmt.Println(a, b)
	}, "/m")

	GET(func() interface{} {
		f, _ := os.Open("C:/Users/Administrator/api/README.MD")
		return NewFileStream(f)
	}, "/download")

	StartService(":8011")
}

func TestParam(t *testing.T) {
	GET(hello1, "/s")
	GET(hello1, "/s/<kk>")
	POST(MulFile, "/update")
	StartService(":8080")
	//reflect.ValueOf(a.show)
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
	type kk int
	fmt.Println(reflect.TypeOf((*kk)(nil)).Elem().Kind())
}
