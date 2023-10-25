package def

type HandlerOrder uint

const (
	Handler_API    = HandlerOrder(100)
	Handler_STATIC = HandlerOrder(99)
)
