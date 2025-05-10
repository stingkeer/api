package rest

import (
	"net/http"
	"testing"

	"go.aew.app/api.v1"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/test/r"
)

func TestNewRedirect(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			return api.NewRedirect("https://www.google.com")
		}, "/redirect")
	}).DoRequestNobody(func(resp *r.Response) {
		if resp.Code() != http.StatusFound && resp.Header("host") != "https://www.google.com" {
			t.Errorf("expect code 302 but %d", resp.Code())
		}
	})

}
