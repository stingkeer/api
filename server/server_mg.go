package server

import (
	"go.uber.org/dig"
)

var (
	c = dig.New()
)

func Provide(constructor interface{}, opts ...dig.ProvideOption) {
	c.Provide(constructor, opts...)
}

func Invoke(constructor interface{}, opts ...dig.InvokeOption) {
	c.Invoke(constructor, opts...)
}
