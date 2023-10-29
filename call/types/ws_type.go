package types

import (
	"net/http"
	"reflect"

	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/kit/ws"
	"gitee.com/fast_api/api/log"
	"github.com/gorilla/websocket"
)

var _ def.Adapter = (*WSType)(nil)

type WSType struct {
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Mapper implements def.Adapter.
func (*WSType) Mapper(param *def.ParamWarp) reflect.Value {
	if websocket.IsWebSocketUpgrade(param.Request.Request) {
		c, err := upgrader.Upgrade(param.Request.ResponseWriter(), param.Request.Request, nil)
		if err != nil {
			panic(err)
		}
		if param.PTyp.Kind() == reflect.Ptr {
			return reflect.ValueOf(ws.NewWSCtx(nil, c))
		}
		if param.PTyp.Kind() == reflect.Struct {
			return reflect.ValueOf(ws.NewWSCtx(nil, c)).Elem()
		}
		c.SetCloseHandler(func(code int, text string) error {
			log.Debug(code, text)
			return nil
		})

	}
	panic("Not a websocket connection")
}

// Register implements def.Adapter.
func (*WSType) Register() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*ws.WSCtx)(nil)).Elem(),
		reflect.TypeOf((*ws.WSCtx)(nil)),
	}
}
