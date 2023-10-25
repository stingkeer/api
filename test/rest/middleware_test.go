package rest

import (
	"net/http"
	"testing"

	"gitee.com/fast_api/api"
)

func handle1() string {
	return "hello,word"
}

func handle2() {

}

func loginAuth(rw http.ResponseWriter, req *http.Request) bool {
	return true
}

func Test_middleware(t *testing.T) {
	api.AddRoutes(
		api.GET(handle1, "/login"),
		api.GET(handle2, "/showOrder"),
	).Middleware(loginAuth)
}
