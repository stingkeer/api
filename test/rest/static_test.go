package rest

import (
	"net/http"
	"testing"

	"gitee.com/fast_api/api"
)

func TestStatic(t *testing.T) {
	api.Static("/web/*", "web", http.Dir("."))
	api.StartService(nil)
}
