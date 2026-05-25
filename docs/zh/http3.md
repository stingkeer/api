# HTTP/3 (QUIC)

框架通过 `quic-go` 库支持基于 QUIC 协议的 HTTP/3。

## 构建 HTTP/3 版本

```bash
go build -tags "http3"
```

这会编译 `http3.go` 替代 `http.go`，使用双 HTTP/1.1 + HTTP/3 服务器替代标准 HTTP 服务器。

## 启动服务

HTTP/3 模式下仅 `StartTLSService` 可用 — QUIC 需要 TLS：

```go
api.StartTLSService(
    api.WithListen(":443"),
    api.WithTLSFile("cert.pem", "key.pem"),
)
```

`StartService()` 在 HTTP/3 模式下返回错误：`"http3 only support tls"`。

## 工作原理

使用 `http3` 构建标签时，`StartTLSService` 同时启动**两个服务器**：

1. **TCP + TLS 服务器** — 标准的 HTTPS，在配置的地址上监听
2. **QUIC 服务器** — HTTP/3，在相同地址上使用 UDP 监听

TCP 服务器自动添加 `Alt-Svc` 响应头宣告 HTTP/3 可用：

```
Alt-Svc: h3=":443"; ma=2592000
```

告知浏览器后续请求可以使用 HTTP/3。

## 错误处理

框架等待任一服务器失败：

```go
select {
case err := <-hErr:
    quicServer.Close()
    return err
case err := <-qErr:
    return err
}
```

任一服务器失败时，另一个会被关闭。

## 客户端支持

大多数现代浏览器支持 HTTP/3：
- Chrome（自 v87 起）
- Firefox（自 v88 起）
- Safari（自 v16 起）
- Edge（自 v87 起）

`curl` 使用 `--http3` 标志支持 HTTP/3：

```bash
curl --http3 https://localhost:443/api/data
```
