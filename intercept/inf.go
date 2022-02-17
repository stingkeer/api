package intercept

import "net/http"

type HttpIntercept interface {
	// Http /**
	Http(rw http.ResponseWriter, req *http.Request) bool
	Order() int
}

type MethodIntercept interface {
	BeforeInvoke()
	AfterInvoke()
}
