package api

import "testing"

func hello() string {
	return "sssss"
}

func TestBind(t *testing.T) {
	GetApi().GET(hello, "/hello")
	Start("")
}
