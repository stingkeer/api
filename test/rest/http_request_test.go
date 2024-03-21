package rest

import (
	"fmt"
	"net/http"
	"testing"

	"gitee.com/fast_api/api"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/test/r"
)

func TestRequest(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(req http.Request) any {
			return req.URL
		}, "/request")
	}).Request().AddHeader("token", "asdfed").Do(func(resp *r.Response) {
		fmt.Println(resp.BodyString())
	})
}

func TestRequestPtr(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(req *http.Request) any {
			return req.URL
		}, "/request")
	}).Request().AddHeader("token", "asdfed").Do(func(resp *r.Response) {
		fmt.Println(resp.BodyString())
	})
}
