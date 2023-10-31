package rest

import (
	"testing"

	"gitee.com/fast_api/api"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/test/r"
)

func TestNewRedirect(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			return api.NewRedirect("https://www.google.com")
		}, "/redirect")
	}).DoRequestNobody(func(resp *r.Response) {
		if resp.Code() != 302 && resp.Header("host") != "https://www.google.com" {
			t.Errorf("expect code 302 but %d", resp.Code())
		}
	})

}
