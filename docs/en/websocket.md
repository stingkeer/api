# WebSocket

The framework provides built-in WebSocket support with automatic protocol upgrade.

## Basic Usage

Use `ws.WSCtx` as a function parameter to automatically upgrade the connection:

```go
import "go.aew.app/api.v1/kit/ws"

api.GET(func(w *ws.WSCtx) {
    go w.Receive(func(messageType int, p []byte) {
        log.Println("received:", string(p))
    })

    w.WriteJSON(map[string]string{"message": "hello"})
}, "/ws")
```

Connect from JavaScript:

```javascript
const ws = new WebSocket("ws://localhost:8080/ws");
ws.onmessage = (event) => console.log(event.data);
ws.send("hello");
```

## WSCtx API

| Method | Description |
|--------|-------------|
| `Receive(f func(int, []byte))` | Blocking message receive loop |
| `Send(o any) error` | Send using default serialization |
| `WriteJSON(o any)` | Send as JSON message |
| `Read(p []byte) (int, error)` | Read binary message (io.Reader) |
| `Write(p []byte) (int, error)` | Write binary message (io.Writer) |
| `SetSerialize(s def.Serialize)` | Override serializer for Send() |
| `SetWsLabel(label string)` | Register to connection pool |

## Receiving Messages

`Receive()` is a blocking loop that calls the callback for each message:

```go
api.GET(func(w *ws.WSCtx) {
    go w.Receive(func(messageType int, p []byte) {
        switch messageType {
        case websocket.TextMessage:
            log.Println("text:", string(p))
        case websocket.BinaryMessage:
            log.Println("binary:", len(p), "bytes")
        }
    })

    // Handler returns, but Receive continues in goroutine
}, "/ws")
```

The loop exits when the connection closes.

## Sending Messages

### JSON

```go
w.WriteJSON(map[string]string{"status": "ok"})
```

### Binary (using Serialize)

```go
w.Send(map[string]string{"data": "value"})
// Serializes using the configured Serialize (default: JSON)
```

### Raw Binary

```go
w.Write([]byte{0x01, 0x02, 0x03})
```

## Connection Pool

Register WebSocket connections with labels and retrieve them later from anywhere:

```go
// Register
api.GET(func(w *ws.WSCtx) {
    w.SetWsLabel("user_123")
    go w.Receive(func(mt int, p []byte) {
        // handle messages
    })
}, "/ws")

// Retrieve elsewhere and send
ctx := ws.GetCtx("user_123")
if ctx != nil {
    ctx.Send(map[string]string{"notification": "New message"})
}
```

## Connection Lifecycle

- **Pong handler** — automatically extends read deadline by 60s on each pong
- **Close handler** — removes the connection from the pool on close
- The `gorilla/websocket.Upgrader` allows all origins (`CheckOrigin` returns true)

## Using with Other Parameters

WebSocket parameters can be combined with other parameters:

```go
api.GET(func(token string, w *ws.WSCtx) {
    // token is validated before WebSocket upgrade
    if token == "" {
        return // handler returns, WSCtx not used
    }
    go w.Receive(func(mt int, p []byte) {
        // handle
    })
}, "/ws")
```

The framework's `WsCaller` checks for WebSocket upgrade and handles it transparently.
