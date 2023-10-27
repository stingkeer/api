package rest

import (
	"fmt"
	"testing"

	"gitee.com/fast_api/api"
	"gitee.com/fast_api/api/kit/ws"
	r "gitee.com/fast_api/api/test/R"
)

type People struct {
}

func TestWS(t *testing.T) {
	r.Test(func() {
		api.GET(func(ws *ws.WSCtx) {
			ws.Receive(func(messageType int, p []byte) {
				fmt.Println(messageType, string(p))
			})
		}, "/ws")
	})
}
