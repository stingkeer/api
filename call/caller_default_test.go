package call

import (
	"fmt"
	"gitee.com/fast_api/api/convert"
	"gitee.com/fast_api/api/public"
	"math/big"
	"reflect"
	"testing"
)

type Person struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

func TestCaller(t *testing.T) {
	c := NewCaller(&convert.JsonConvertImpl{}, &convert.DefaultTypeConvert{})
	c.callPost(&public.Entry{
		Url:    "",
		Group:  "",
		Method: "",
		Fn:     show,
		Ids:    nil,
	}, nil)

}

func TestTypeConvert(t *testing.T) {
	c := NewCaller(&convert.JsonConvertImpl{}, &convert.DefaultTypeConvert{})
	v := c.paramTypeConvert("10000000000000", reflect.TypeOf(new(big.Int)))
	fmt.Println(v)
	dd := c.paramTypeConvert("100", reflect.TypeOf(new(int)))
	fmt.Println(dd.Elem())
	cc := c.paramTypeConvert("api", reflect.TypeOf(new(string)).Elem())
	fmt.Println(cc)
}

func show(p *Person) string {
	fmt.Println(p)
	return "shang"
}
