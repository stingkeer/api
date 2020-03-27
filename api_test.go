package api

import "testing"

func hello() string {
	return "shangzebei"
}

func TestBind(t *testing.T) {
	GetApi().GET(hello, "/hello")
	GetApi().PUT(hello, "/hello")
	Start("")

}
