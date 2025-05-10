package rest

import (
	"net/http"
	"testing"

	"go.aew.app/api.v1"
)

func TestStatic(t *testing.T) {
	api.Static("/web/*", "web", http.Dir("."))
	api.StartService(nil)
}
