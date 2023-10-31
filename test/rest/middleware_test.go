package rest

import (
	"fmt"
	"net/http"
	"testing"

	"gitee.com/fast_api/api"
)

func handle1() string {
	return "hello,word"
}

func handle2() string {
	return "handle2"
}

func LoginAuth(req *http.Request) (ret any) {
	fmt.Println("LoginAuth")
	return nil
}

func LoginCookie(req *http.Request) (ret any) {
	fmt.Println("LoginCookie")
	return "Cookie not find"
}

func Test_middleware(t *testing.T) {
	api.AddRoutes(
		api.GET(handle1, "/login"),
		api.GET(handle2, "/showOrder"),
	).Middleware(LoginAuth, LoginCookie)
	api.StartService()
}
