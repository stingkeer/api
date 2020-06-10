package api

import "testing"

func hello() interface{} {
	return map[string]string{"name": "shangzebei"}
}

func hello1(kk string) interface{} {
	return map[string]string{"name": kk}
}

func TestBind(t *testing.T) {
	GetApi().GET(hello, "hello")
	Start(":8080")
}

func TestParam(t *testing.T) {
	GetApi().GET(hello1, "/s")
	Start(":8080")
}
