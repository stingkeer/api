package api

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"testing"
)

func hello() interface{} {
	return map[string]string{"name": "shangzebei"}
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
	GetApi().GET(hello, "/hello")
	Start(":8080")
}

func TestParam(t *testing.T) {
	GetApi().GET(hello1, "/s")
	GetApi().GET(hello1, "/s/<kk>")
	GetApi().POST(MulFile, "/update")
	Start(":8080")
}

func TestFile(t *testing.T) {
	GetApi().POST(MulFile, "/update")
	Start(":8080")
}
