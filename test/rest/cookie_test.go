package rest

import (
	"testing"

	"go.aew.app/api"
	"go.aew.app/api/def"
	"go.aew.app/api/test/r"
)

func TestCookie(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(header def.Header) {
			cookie, err := header.Cookie("username")
			if err != nil {
				return
			}
			cookie.Value = "hello"
			header.SetCookie(cookie)
		}, "/cookie")
	}).Request().SetCookie("username", "my").Do(func(resp *r.Response) {
		cookie := resp.Cookies()[0]
		if cookie.Name != "username" || cookie.Value != "hello" {
			t.Error("TestCookie Error")
		}
	})
}
