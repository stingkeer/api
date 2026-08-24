package rest

import (
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"go.aew.app/api.v1"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/kit/ws"
	"go.aew.app/api.v1/test/r"
)

// dialWs connects a real websocket client to the embedded test server.
func dialWs(t *testing.T, path string) *websocket.Conn {
	t.Helper()
	c, _, err := websocket.DefaultDialer.Dial(r.WSURL(path), nil)
	if err != nil {
		t.Fatalf("dial %s: %v", path, err)
	}
	return c
}

func TestWS(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(w *ws.WSCtx) {
			w.Receive(func(messageType int, p []byte) {
				_ = w.Send(string(p))
			})
		}, "/ws-echo")
	})

	c := dialWs(t, "/ws-echo")
	defer c.Close()

	if err := c.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
		t.Fatal(err)
	}
	_ = c.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, msg, err := c.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	// Plain strings are serialized as text/plain, so the echo comes back
	// exactly as sent (no JSON quoting).
	if got := strings.TrimSpace(string(msg)); got != `ping` {
		t.Errorf("echo: expect %q but %q", `ping`, got)
	}
}

func TestNOPendingWS(t *testing.T) {
	registered := make(chan struct{})
	r.Test(t, func() def.Option {
		return api.GET(func(w *ws.WSCtx) {
			if err := w.SetWsLabel("id"); err != nil {
				t.Errorf("SetWsLabel: %v", err)
			}
			close(registered)
		}, "/ws-idle")
	})

	c := dialWs(t, "/ws-idle")
	defer c.Close()

	select {
	case <-registered:
	case <-time.After(5 * time.Second):
		t.Fatal("handler did not run within 5s")
	}
	if ws.GetCtx("id") == nil {
		t.Error("GetCtx(\"id\") = nil after SetWsLabel")
	}
}

func TestBlockIO(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(w *ws.WSCtx) {
			w.Receive(func(messageType int, p []byte) {
				if string(p) == "hello" {
					_ = w.Send(map[string]string{"aaa": "bbbb"})
				}
			})
		}, "/ws-block")
	})

	c := dialWs(t, "/ws-block")
	defer c.Close()

	if err := c.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
		t.Fatal(err)
	}
	_ = c.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, msg, err := c.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(msg)); got != `{"aaa":"bbbb"}` {
		t.Errorf("block io: expect %s but %s", `{"aaa":"bbbb"}`, got)
	}
}

func TestPanic(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(w *ws.WSCtx) {
			panic("ws error")
		}, "/ws-panic")
	})

	c := dialWs(t, "/ws-panic")
	defer c.Close()

	// A panic inside the handler must not hang the connection: the client
	// should observe the connection close (read error) in bounded time.
	_ = c.SetReadDeadline(time.Now().Add(5 * time.Second))
	if _, _, err := c.ReadMessage(); err == nil {
		t.Error("expect connection closed after handler panic, but read succeeded")
	}
}
