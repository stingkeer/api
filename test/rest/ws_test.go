package rest

import (
	"fmt"
	"testing"
	"time"

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

func TestNOPendingWS(t *testing.T) {
	go func() {
		time.Sleep(time.Second * 30)
		ws.GetCtx("id").Send("sadfasdf")
	}()
	r.Test(func() {
		api.GET(func(ws *ws.WSCtx) {
			ws.SetWsLabel("id")
		}, "/ws")
	})
}

func TestBlockIO(t *testing.T) {
	r.Test(func() {
		api.GET(func(ws *ws.WSCtx) {
			ws.Receive(func(messageType int, p []byte) {
				if string(p) == "hello" {
					ws.Send(map[string]string{"aaa": "bbbb"})
				}
			})
		}, "/ws")
	})
}

// TODO Fix Me
func TestPanic(t *testing.T) {
	r.Test(func() {
		api.GET(func(ws *ws.WSCtx) {
			panic("ws error")
		}, "/ws")
	})
}
