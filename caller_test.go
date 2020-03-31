package api

import (
	"fmt"
	"testing"
)

type Person struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

func TestCaller(t *testing.T) {

	c := CallerDefault{&JSONConvertImpl{}}
	c.callPost(show, nil)

}

func show(p *Person) string {
	fmt.Println(p)
	return "shang"
}
