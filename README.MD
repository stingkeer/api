[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/stingkeer/api)

# api.v1 - Go Declarative Web Framework

A declarative Go web API framework that infers parameter binding from function signatures using DWARF debug information, with zero runtime reflection overhead for parameter name resolution.

**[English Documentation](docs/en/)** | **[中文文档](docs/zh/)**

---

## Features

- **Declarative Routing** — Function signatures automatically define URL parameter binding
- **DWARF Parameter Resolution** — Reads DWARF debug info from compiled binary for parameter names
- **Radix Tree Router** — High-performance routing with static, parameterized, and regex-constrained paths
- **Multi-Protocol** — HTTP/1.1, HTTPS/TLS, HTTP/3 (QUIC)
- **WebSocket** — Built-in support with automatic protocol upgrade
- **Swagger** — Auto-generated OpenAPI 3.0 documentation
- **Compression** — Automatic gzip/deflate response compression
- **Caching** — Extensible method-level cache with custom persistence

## Quick Start

```bash
go get go.aew.app/api.v1
```

```go
package main

import "go.aew.app/api.v1"

func main() {
    api.GET(func() interface{} {
        return "hello world"
    }, "/hello")

    api.StartService(api.WithListen(":8080"))
}
```

## Auto Parameter Binding

```go
// Query params — names resolved via DWARF debug info
api.GET(func(username, password string) interface{} {
    return map[string]string{"username": username, "password": password}
}, "/login")

// Path params
api.GET(func(id string) interface{} {
    return map[string]string{"id": id}
}, "/user/<id>")

// Request body
api.POST(func(body LoginRequest) interface{} {
    return body
}, "/login")

// File download with Range support and rate limiting
api.GET(func() interface{} {
    f, _ := os.Open("video.mp4")
    return api.NewStream(f).SetName("video.mp4").SetRateLimit(1024 * 1024)
}, "/download")

// WebSocket
api.GET(func(w *ws.WSCtx) {
    go w.Receive(func(mt int, p []byte) { log.Println(string(p)) })
    w.WriteJSON(map[string]string{"message": "hello"})
}, "/ws")
```

## Server Configuration

```go
// HTTP
api.StartService(api.WithListen(":8080"))

// HTTPS
api.StartTLSService(api.WithListen(":443"), api.WithTLSFile("cert.pem", "key.pem"))
```

## Build Tags

```bash
go build -tags 'swagger'   # Enable Swagger UI
go build -tags 'http3'     # Enable HTTP/3 (QUIC)
```

## Documentation

| Topic | English | 中文 |
|-------|---------|------|
| Getting Started | [Quick Start](docs/en/quick-start.md) | [快速开始](docs/zh/quick-start.md) |
| Routing | [Routing](docs/en/routing.md) | [路由注册](docs/zh/routing.md) |
| Parameters | [Parameter Binding](docs/en/parameter-binding.md) | [参数绑定](docs/zh/parameter-binding.md) |
| Return Types | [Return Types](docs/en/return-types.md) | [返回类型](docs/zh/return-types.md) |
| Middleware | [Middleware](docs/en/middleware.md) | [中间件](docs/zh/middleware.md) |
| Static Files | [Static Files](docs/en/static-files.md) | [静态文件服务](docs/zh/static-files.md) |
| WebSocket | [WebSocket](docs/en/websocket.md) | [WebSocket](docs/zh/websocket.md) |
| Upload & Download | [File Upload & Download](docs/en/file-upload-download.md) | [文件上传与下载](docs/zh/file-upload-download.md) |
| Caching | [Caching](docs/en/caching.md) | [缓存系统](docs/zh/caching.md) |
| Configuration | [Server Configuration](docs/en/server-configuration.md) | [服务启动配置](docs/zh/server-configuration.md) |
| HTTP/3 | [HTTP/3 (QUIC)](docs/en/http3.md) | [HTTP/3 (QUIC)](docs/zh/http3.md) |
| Swagger | [Swagger](docs/en/swagger.md) | [Swagger 文档](docs/zh/swagger.md) |
| Errors | [Error Handling](docs/en/error-handling.md) | [错误处理](docs/zh/error-handling.md) |
| Compression | [Compression](docs/en/compression.md) | [响应压缩](docs/zh/compression.md) |
| Logging | [Logging](docs/en/logging.md) | [日志系统](docs/zh/logging.md) |
| Architecture | [Architecture](docs/en/architecture.md) | [架构设计](docs/zh/architecture.md) |

## License

See [LICENSE](LICENSE) for details.
