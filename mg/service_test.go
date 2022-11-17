package mg

import (
	"fmt"
	"testing"
)

type A struct {
	Name string
}

func TestMg(t *testing.T) {
	err := Provide(func() *A {
		return &A{Name: "service"}
	})
	if err != nil {
		panic(err)
	}

	err = Invoke(func(a *A) {
		fmt.Println("#####", a.Name)
	})
	if err != nil {
		panic(err)
	}
}
