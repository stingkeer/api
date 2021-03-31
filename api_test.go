package api

import (
	"fmt"
	"gitee.com/fast_api/api/public"
	"gitee.com/fast_api/api/server"
	"github.com/sirupsen/logrus"

	"io/ioutil"
	"mime/multipart"

	"testing"
)

func hello() interface{} {
	//panic("asdfsadf")
	return "asdfsd"
}

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

func TestBind(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	GET(hello, "/hello")
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

func (A) Encode(interface{}) *public.Content {
	return nil
}
func (A) Decode([]byte, interface{}) error {
	return nil
}

func TestFile(t *testing.T) {
	server.Provide(func() public.Serialize {
		return &A{}
	})
	POST(MulFile, "/update")
	StartService(":8080")
}

func TestFx(t *testing.T) {

}
