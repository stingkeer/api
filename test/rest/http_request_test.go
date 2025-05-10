package rest

import (
	"fmt"
	"net/http"
	"testing"

	"go.aew.app/api.v1"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/test/r"
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
