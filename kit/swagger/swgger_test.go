package swagger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestDataType(t *testing.T) {
	var params []ParameterObject = []ParameterObject{}
	fmt.Println(params)
	dd, _ := json.Marshal(params)
	fmt.Println(string(dd))
}

func TestDefine(t *testing.T) {
	var a = struct {
		Name string
	}{}
	definitions(reflect.TypeOf(a))
}

func TestStructParam(t *testing.T) {
	type S struct {
		Page int `json:"page"`
		Size int `json:"size"`
	}

	type B struct {
		S
		A int `json:"a"`
		B int `json:"b"`
	}
	xx := genStructParam(reflect.TypeOf(B{}))
	fmt.Println(xx)
}
