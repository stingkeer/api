package http

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestRef(t *testing.T) {
	fmt.Println(reflect.TypeOf(errors.New("sadfsdf")).Elem().String())
}
