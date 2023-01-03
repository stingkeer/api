package mg

import (
	"go.uber.org/dig"
)

var (
	mg      ServiceMg
	Provide = func(constructor interface{}) error { return mg.Provide(constructor) }
	Invoke  = func(constructor interface{}) error { return mg.Invoke(constructor) }
)

func init() {
	mg = &MgImpl{c: dig.New()}
}

type MgImpl struct {
	c *dig.Container
}

func (m *MgImpl) Provide(constructor interface{}) error {
	return m.c.Provide(constructor)
}

func (m *MgImpl) Invoke(constructor interface{}) error {
	return m.c.Invoke(constructor)
}
