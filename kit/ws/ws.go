package ws

import (
	"errors"
	"io"
	"time"

	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/log"
	"github.com/gorilla/websocket"
)

var (
	_        io.ReadWriter = (*WSCtx)(nil)
	pongWait               = 60 * time.Second
)

type WSCtx struct {
	serialize def.Serialize
	conn      *websocket.Conn
	label     string
}

func NewWSCtx(conn *websocket.Conn) *WSCtx {
	x := &WSCtx{serialize: def.DefaultContext.Serialize, conn: conn}
	x.init()
	return x
}

// Read implements io.ReadWriter. only support BinaryMessage
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

// Write implements io.ReadWriter. only support BinaryMessage
func (ws *WSCtx) Write(p []byte) (n int, err error) {
	writer, err := ws.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return 0, err
	}
	return writer.Write(p)
}

func (ws *WSCtx) SetSerialize(serialize def.Serialize) {
	ws.serialize = serialize
}

// Receive This method will clog up
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

func (ws *WSCtx) SetWsLabel(label string) *WSCtx {
	setWs(label, ws)
	ws.label = label
	return ws
}

func (ws *WSCtx) init() {

	ws.conn.SetPongHandler(func(string) error {
		ws.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ws.conn.SetCloseHandler(func(code int, text string) error {
		log.Debug(code, text)
		delete(cPool, ws.label)
		return nil
	})

}
