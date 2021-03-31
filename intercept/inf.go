package intercept

import "net/http"

type HttpIntercept interface {
	/**
	type return bool is true
	*/
	Http(rw http.ResponseWriter, req *http.Request) bool
	Order() int
}

type MethodIntercept interface {
}
