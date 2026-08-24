package ws

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/log"
)

var (
	_          io.ReadWriter = (*WSCtx)(nil)
	pongWait                 = 60 * time.Second
	writeWait                = 10 * time.Second
	pingPeriod               = (pongWait * 9) / 10

	// DefaultReadLimit caps inbound message size (DoS guard: without a limit a
	// client can buffer unbounded memory server-side). Override per connection
	// with SetReadLimit, or package-wide by assigning this var.
	DefaultReadLimit int64 = 1 << 20 // 1MB
)

type WSCtx struct {
	serialize def.Serialize
	conn      *websocket.Conn
	label     string
	done      chan struct{}
	once      sync.Once
	writeMu   sync.Mutex
	reader    io.Reader
	onError   func(err error)
}

func NewWSCtx(conn *websocket.Conn) *WSCtx {
	x := &WSCtx{
		serialize: def.DefaultContext.Serialize,
		conn:      conn,
		done:      make(chan struct{}),
	}
	x.init()
	return x
}

func (ws *WSCtx) SetOnError(f func(err error)) {
	ws.onError = f
}

func (ws *WSCtx) init() {
	ws.conn.SetReadLimit(DefaultReadLimit)
	ws.conn.SetReadDeadline(time.Now().Add(pongWait))

	ws.conn.SetPongHandler(func(string) error {
		ws.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ws.conn.SetCloseHandler(func(code int, text string) error {
		log.Debugf("websocket close: code=%d text=%s", code, text)
		ws.writeMu.Lock()
		ws.conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(code, ""),
			time.Now().Add(writeWait))
		ws.writeMu.Unlock()
		ws.Close()
		return nil
	})

	go ws.pingLoop()
}

// SetReadLimit overrides the inbound message size cap for this connection.
func (ws *WSCtx) SetReadLimit(n int64) { ws.conn.SetReadLimit(n) }

func (ws *WSCtx) pingLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ws.writeMu.Lock()
			err := ws.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait))
			ws.writeMu.Unlock()
			if err != nil {
				// Peer is gone: tear the connection down (deadline cleanup,
				// pool deregistration) instead of leaking until the read
				// deadline fires.
				ws.Close()
				return
			}
		case <-ws.done:
			return
		}
	}
}

func (ws *WSCtx) Read(p []byte) (n int, err error) {
	if ws.reader == nil {
		_, reader, err := ws.conn.NextReader()
		if err != nil {
			return 0, err
		}
		ws.reader = reader
	}
	n, err = ws.reader.Read(p)
	if err != nil {
		// readers may wrap EOF; compare via errors.Is
		if errors.Is(err, io.EOF) {
			ws.reader = nil
		}
		return n, err
	}
	return n, nil
}

func (ws *WSCtx) Write(p []byte) (n int, err error) {
	ws.writeMu.Lock()
	defer ws.writeMu.Unlock()
	writer, err := ws.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return 0, err
	}
	n, err = writer.Write(p)
	if closeErr := writer.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	return
}

func (ws *WSCtx) SetSerialize(serialize def.Serialize) {
	ws.serialize = serialize
}

func (ws *WSCtx) Receive(f func(messageType int, p []byte)) {
	defer ws.Close()
	for {
		mt, message, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Errorf("websocket read error: %v", err)
				if ws.onError != nil {
					ws.onError(err)
				}
			}
			return
		}
		f(mt, message)
	}
}

func (ws *WSCtx) Send(o any) error {
	content := ws.serialize.Encode(o)
	if content == nil {
		return errors.New("encode returned nil content")
	}
	msgType := websocket.TextMessage
	if !strings.Contains(content.ContentType, "json") && !strings.Contains(content.ContentType, "text") {
		msgType = websocket.BinaryMessage
	}
	ws.writeMu.Lock()
	defer ws.writeMu.Unlock()
	return ws.conn.WriteMessage(msgType, content.Bytes)
}

func (ws *WSCtx) WriteJSON(o any) error {
	ws.writeMu.Lock()
	defer ws.writeMu.Unlock()
	return ws.conn.WriteJSON(o)
}

func (ws *WSCtx) SetWsLabel(label string) error {
	if !setWs(label, ws) {
		return fmt.Errorf("websocket: failed to register label %q (duplicate or max connections reached)", label)
	}
	ws.label = label
	return nil
}

func (ws *WSCtx) Close() error {
	var err error
	ws.once.Do(func() {
		ws.writeMu.Lock()
		err = ws.conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			time.Now().Add(writeWait))
		ws.writeMu.Unlock()
		close(ws.done)
		ws.conn.Close()
		if ws.label != "" {
			pool.delete(ws.label)
		}
	})
	return err
}
