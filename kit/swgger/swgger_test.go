package swgger

import (
	"reflect"
	"testing"
)

func TestDataType(t *testing.T) {

}

func TestDefine(t *testing.T) {
	var a = struct {
		Name string
	}{}
	definitions(reflect.TypeOf(a))
}
