package types

import (
	"net/http"
	"reflect"

	"github.com/gorilla/websocket"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/kit/ws"
	"go.aew.app/api.v1/log"
)

var _ def.Adapter = (*WSType)(nil)

type WSType struct {
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (*WSType) Mapper(param *def.ParamWarp) reflect.Value {
	if !websocket.IsWebSocketUpgrade(param.Request.Request) {
		log.Errorf("websocket: not a websocket connection for %s", param.Request.URL)
		return reflect.Value{}
	}
	c, err := upgrader.Upgrade(param.Request.ResponseWriter(), param.Request.Request, nil)
	if err != nil {
		log.Errorf("websocket upgrade failed: %v", err)
		return reflect.Value{}
	}
	if param.PTyp.Kind() == reflect.Ptr {
		return reflect.ValueOf(ws.NewWSCtx(c))
	}
	if param.PTyp.Kind() == reflect.Struct {
		return reflect.ValueOf(ws.NewWSCtx(c)).Elem()
	}
	log.Errorf("websocket: unsupported parameter type %v", param.PTyp)
	return reflect.Value{}
}

func (*WSType) Register() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*ws.WSCtx)(nil)).Elem(),
		reflect.TypeOf((*ws.WSCtx)(nil)),
	}
}
