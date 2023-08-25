package api

import (
	"crypto/tls"
	"fmt"
	"gitee.com/fast_api/api/cache"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/mg"
	"github.com/sirupsen/logrus"
	"io"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

func hello1(kk def.StringReq) any {
	return map[string]string{"name": kk.String()}
}

func MulFile(read multipart.Reader) string {
	par, _ := read.NextPart()
	fmt.Println(par.FileName(), par.FormName())
	b, _ := io.ReadAll(par)
	fmt.Println(string(b))
	return string(b)
}

func TestApiHttp(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	GET(func(a http.Request) (any, any) {
		return a.Host, cache.NewCacheImpl(time.Second)
	}, "/hello")

	GET(func(a big.Int) any {
		return a.String()
	}, "/int")

	GET(func(a def.Header) any {
		return a.Values("Accept-Encoding")
	}, "/h")

	GET(func(a def.Header) {
		a.Add("szb", "nnnnn")
	}, "/no")

	GET(func(a int, hello def.String[cache.Key]) any {
		fmt.Println(a, hello, hello.String())
		return nil
	}, "/m")

	GET(func(req *http.Request, resp http.Response) {

	}, "/http")

	GET(func() any {
		f, e := os.Open("d:/download/QmfWv8FfpKiCWsueKfXDLrgyqXZsEuGFJFBL7TfjNmxkAw")
		fmt.Println(e)
		stream := NewStream(f)
		stream.SetRateLimit(500000)
		return stream
	}, "/download")

	GET(func(header def.Header, reader multipart.Reader) {
		fmt.Println(header, reader)
	}, "/file")

	StartService(nil)
}

func TestURL(t *testing.T) {
	GET(hello1, "/s")
	GET(hello1, "/s/<kk>")
	POST(MulFile, "/update")
	StartService(nil)
}

func TestPath(t *testing.T) {
	GET(hello1, "/s/<kk>")
	StartService(func(conf *Config) *Config {
		return conf
	})
}

type A struct {
	_ string
}

func (A) Encode(any) *def.Content {
	return nil
}
func (A) Decode([]byte, any) error {
	return nil
}

func TestFile(t *testing.T) {
	mg.Provide(func() def.Serialize {
		return &A{}
	})
	POST(MulFile, "/update")
	StartService(nil)
}

func TestType(t *testing.T) {
	//bytes 64448694-86052627/86052628
	//21603934
	//13f63f85d5b618b2c03121c0a8d38758
	//13f63f85d5b618b2c03121c0a8d38758
	//2-9
}

func TestDail(t *testing.T) {
	c, err := tls.Dial("tcp", "www.baidu.com:https", nil)
	fmt.Println(c, err, c.VerifyHostname("www.baidu.com"))
}

func TestApiAfter(t *testing.T) {
	http.FileServer(http.Dir(`.`))
	go StartService(nil)
	GET(hello1, "/s")
}

func TestHtml(t *testing.T) {
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{range .Items}}<div>{{ . }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}
	</body>
</html>`
	GET(func() any {
		data := struct {
			Title string
			Items []string
		}{
			Title: "My page",
			Items: []string{
				"My photos",
				"My blog",
			},
		}
		return Html(tpl, data)
	}, "/html")
	StartService(nil)
}

func TestStatic(t *testing.T) {
	Static("/web/*", "web", http.Dir("."))
	StartService(nil)
}

func TestNewRedirect(t *testing.T) {
	GET(func() any {
		return NewRedirect("https://www.google.com")
	}, "/redirect")
	StartService(nil)
}

func TestCookie(t *testing.T) {
	//cookie
	GET(func(header def.Header) {
		cookie, err := header.Cookie("username")
		if err != nil {
			return
		}
		cookie.Value = "hello"
		header.SetCookie(cookie)
	}, "/cookie")
	StartService(nil)
}

func TestCache(t *testing.T) {
	GET(func(s def.String[cache.Key]) (any, cache.Cache) {
		fmt.Println("invoke")
		return "hello", cache.NewCacheImpl(time.Second * 30)
	}, "/cache")
	StartService(nil)
}

func TestAfterRegister(t *testing.T) {

	var s sync.WaitGroup
	s.Add(1)
	go StartService(nil)
	GET(func(s def.StringReq) any {
		return "after" + s.String()
	}, "/after")
	s.Wait()

}
