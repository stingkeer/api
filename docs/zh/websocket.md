# WebSocket

框架提供内建的 WebSocket 支持，自动处理协议升级。

## 基本用法

使用 `ws.WSCtx` 作为函数参数即可自动升级连接：

```go
import "go.aew.app/api.v1/kit/ws"

api.GET(func(w *ws.WSCtx) {
    go w.Receive(func(messageType int, p []byte) {
        log.Println("received:", string(p))
    })

    w.WriteJSON(map[string]string{"message": "hello"})
}, "/ws")
```

JavaScript 连接：

```javascript
const ws = new WebSocket("ws://localhost:8080/ws");
ws.onmessage = (event) => console.log(event.data);
ws.send("hello");
```

## WSCtx API

| 方法 | 说明 |
|------|------|
| `Receive(f func(int, []byte))` | 阻塞式消息接收循环 |
| `Send(o any) error` | 使用默认序列化发送 |
| `WriteJSON(o any)` | 发送 JSON 消息 |
| `Read(p []byte) (int, error)` | 读取二进制消息（io.Reader） |
| `Write(p []byte) (int, error)` | 写入二进制消息（io.Writer） |
| `SetSerialize(s def.Serialize)` | 覆盖 Send() 的序列化器 |
| `SetWsLabel(label string)` | 注册到连接池 |

## 接收消息

`Receive()` 是阻塞循环，每条消息调用回调函数：

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

    // handler 返回，但 Receive 在 goroutine 中继续运行
}, "/ws")
```

连接关闭时循环退出。

## 发送消息

### JSON

```go
w.WriteJSON(map[string]string{"status": "ok"})
```

### 二进制（使用序列化器）

```go
w.Send(map[string]string{"data": "value"})
// 使用配置的 Serialize 序列化（默认：JSON）
```

### 原始二进制

```go
w.Write([]byte{0x01, 0x02, 0x03})
```

## 连接池

通过标签注册 WebSocket 连接，随后可在任意位置获取：

```go
// 注册
api.GET(func(w *ws.WSCtx) {
    w.SetWsLabel("user_123")
    go w.Receive(func(mt int, p []byte) {
        // 处理消息
    })
}, "/ws")

// 在其他地方获取并发送
ctx := ws.GetCtx("user_123")
if ctx != nil {
    ctx.Send(map[string]string{"notification": "New message"})
}
```

## 连接生命周期

- **Pong 处理器** — 每次收到 pong 自动延长读超时 60 秒
- **关闭处理器** — 连接关闭时从池中移除
- `gorilla/websocket.Upgrader` 允许所有来源（`CheckOrigin` 返回 true）

## 与其他参数组合使用

WebSocket 参数可以与其他参数组合：

```go
api.GET(func(token string, w *ws.WSCtx) {
    // token 在 WebSocket 升级前验证
    if token == "" {
        return // handler 返回，WSCtx 不使用
    }
    go w.Receive(func(mt int, p []byte) {
        // 处理消息
    })
}, "/ws")
```

框架的 `WsCaller` 会检查 WebSocket 升级并透明处理。
