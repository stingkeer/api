package mg

type ServiceMg interface {
	Provide(constructor interface{}) error
	Invoke(constructor interface{}) error
}
