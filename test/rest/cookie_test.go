package rest

import (
	"testing"

	"gitee.com/fast_api/api"
	"gitee.com/fast_api/api/def"
)

func TestCookie(t *testing.T) {
	//cookie
	api.GET(func(header def.Header) {
		cookie, err := header.Cookie("username")
		if err != nil {
			return
		}
		cookie.Value = "hello"
		header.SetCookie(cookie)
	}, "/cookie")
	api.StartService(nil)
}
