package api

import (
	"fmt"
	"testing"
)

func TestUtils(t *testing.T) {
	s := "main.(*PerSon).show-fm"
	s1 := "main.hello"
	fmt.Println(s, s1)
	fmt.Println(SplitFuncName(s))
}
