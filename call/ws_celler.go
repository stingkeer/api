package call

import (
	"github.com/gorilla/websocket"
	"go.aew.app/api/def"
)

var _ def.Caller = (*WsCaller)(nil)

type WsCaller struct {
	TraceCaller
}

func NewWsCaller(serialize def.Serialize, pool *def.MethodsPools) *WsCaller {
	return &WsCaller{
		TraceCaller: *NewTraceCaller(serialize, pool),
	}
}

// Call implements def.Caller.
func (w *WsCaller) Call(f *def.Entry, req *def.Request) interface{} {
	x := w.TraceCaller.Call(f, req)
	if websocket.IsWebSocketUpgrade(req.Request) {
		return def.Empty("")
	}
	return x
}
