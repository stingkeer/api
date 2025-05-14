package def

type HandlerOrder uint

// HandlerOrder
// 0-100     system
// 100-1000  user
const (
	Handler_API           = HandlerOrder(100)
	Handler_HTTP_COMPRESS = HandlerOrder(1500)
	Handler_STATIC        = HandlerOrder(99)
	Handler_NOTFIND       = HandlerOrder(1000)
)
