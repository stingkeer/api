```go
package main

import (
	"gitee.com/fast_api/api"
	"strconv"
)

func main() {
	api.GetApi().GET(hello, "hello")
	api.GetApi().GET(ttt, "ppp")
	p := &PerSon{}
	api.GetApi().GET(p.show, "person")
	api.GetApi().POST(p.login, "login")
	api.Start(":8080")
}

type Login struct {
	UserName string
	PassWord string
}

type PerSon struct {
}

func hello(show string) string {
	return "shangzebei" + show
}

func pppp() string {
	return "asdfsadfsdfsdfsd"
}

func ttt() map[string]string {
	return map[string]string{"asdfsdf": "sdfsdf"}
}

func (p *PerSon) show(ss int, _pp string) string {
	return "person" + strconv.Itoa(ss)
}

func (p *PerSon) login(l Login) *Login {
	return &l
}
```