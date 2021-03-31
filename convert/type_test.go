package convert

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestType(t *testing.T) {
	v := reflect.ValueOf(big.NewInt(100)).Elem()
	fmt.Println(v.Addr())
}
