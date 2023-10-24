package rest

import (
	"testing"

	"gitee.com/fast_api/api"
)

func TestNewRedirect(t *testing.T) {
	api.GET(func() any {
		return api.NewRedirect("https://www.google.com")
	}, "/redirect")
	api.StartService(nil)
}
