package rest

import (
	"fmt"
	"testing"
	"time"

	"go.aew.app/api"
	"go.aew.app/api/def"
	"go.aew.app/api/kit/ws"
	"go.aew.app/api/test/r"
)

type People struct {
}

func TestWS(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(ws *ws.WSCtx) {
			ws.Receive(func(messageType int, p []byte) {
				fmt.Println(messageType, string(p))
			})
		}, "/ws")
	}).Wait()
}

func TestNOPendingWS(t *testing.T) {
	go func() {
		time.Sleep(time.Second * 30)
		ws.GetCtx("id").Send("sadfasdf")
	}()
	r.Test(t, func() def.Option {
		return api.GET(func(ws *ws.WSCtx) {
			ws.SetWsLabel("id")
		}, "/ws")
	}).Wait()
}

func TestBlockIO(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(ws *ws.WSCtx) {
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
	r.Test(t, func() def.Option {
		return api.GET(func(ws *ws.WSCtx) {
			panic("ws error")
		}, "/ws")
	})
}
