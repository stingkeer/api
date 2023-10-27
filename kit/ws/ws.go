package ws

import (
	"errors"
	"io"

	"gitee.com/fast_api/api/def"
	"github.com/gorilla/websocket"
)

var _ io.ReadWriter = (*WSCtx)(nil)

type WSCtx struct {
	serialize def.Serialize
	conn      *websocket.Conn
}

func NewWSCtx(serialize def.Serialize, conn *websocket.Conn) *WSCtx {
	return &WSCtx{serialize: serialize, conn: conn}
}

// Read implements io.ReadWriter.
func (ws *WSCtx) Read(p []byte) (n int, err error) {
	typ, reader, err := ws.conn.NextReader()
	if err != nil {
		return 0, err
	}
	if typ == websocket.BinaryMessage {
		return reader.Read(p)
	}
	return 0, errors.New("websocket type err")
}

// Write implements io.ReadWriter.
func (ws *WSCtx) Write(p []byte) (n int, err error) {
	writer, err := ws.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return 0, err
	}
	return writer.Write(p)
}

func (ws *WSCtx) Receive(f func(messageType int, p []byte)) {
	for {
		mt, message, err := ws.conn.ReadMessage()
		if err != nil {
			break
		}
		f(mt, message)
	}
}

// Send Use the default serialize send
func (ws *WSCtx) Send(o any) {
	context := ws.serialize.Encode(o)
	ws.conn.WriteMessage(websocket.BinaryMessage, context.Bytes)
}

func (ws *WSCtx) WriteJSON(o any) {
	ws.conn.WriteJSON(o)
}
