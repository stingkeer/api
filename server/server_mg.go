package server

import (
	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/match"
	"gitee.com/fast_api/api/public"
	"gitee.com/fast_api/api/serialize"
	"go.uber.org/dig"
)

var (
	c *dig.Container
)

func init() {
	c = dig.New()
	//default
	Provide(func() public.Serialize {
		return &serialize.JsonConvertImpl{}
	})

	Provide(func(resultConvert public.Serialize) public.Caller {
		return call.NewCaller(resultConvert)
	})

	Provide(func() public.MetaMethods {
		return public.MethodsPools
	})

	Provide(func() public.Match {
		return match.NewMatchImpl()
	})

	public.MethodsPools = make(public.MetaMethods)
}

func Provide(constructor interface{}, opts ...dig.ProvideOption) {
	c.Provide(constructor, opts...)
}

func Invoke(constructor interface{}, opts ...dig.InvokeOption) {
	c.Invoke(constructor, opts...)
}
