package api

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func TestDef(t *testing.T) {
	PackApi()
}

func aa() {

}

func show(f interface{}) {
	v := reflect.ValueOf(f)
	fmt.Println(runtime.FuncForPC(v.Pointer()).Name())
}

func TestFunc(t *testing.T) {
	go func() {
		show(aa)
	}()
}
