# api.v1 - Go Declarative Web Framework

A Go declarative web API framework that automatically infers parameter binding from function signatures using DWARF debug information, with zero runtime reflection overhead for parameter name resolution.

## Table of Contents

1. [Quick Start](quick-start.md)
2. [Routing](routing.md)
3. [Parameter Binding](parameter-binding.md)
4. [Return Types](return-types.md)
5. [Middleware](middleware.md)
6. [Static Files](static-files.md)
7. [WebSocket](websocket.md)
8. [File Upload & Download](file-upload-download.md)
9. [Caching](caching.md)
10. [Server Configuration](server-configuration.md)
11. [HTTP/3 (QUIC)](http3.md)
12. [Swagger Documentation](swagger.md)
13. [Error Handling](error-handling.md)
14. [Compression](compression.md)
15. [Logging](logging.md)
16. [Architecture](architecture.md)

## Core Features

- **Declarative Routing** — Function signatures automatically define URL parameter binding
- **DWARF Parameter Resolution** — Reads DWARF debug info from compiled binary to obtain parameter names, avoiding runtime reflection
- **Radix Tree Router** — High-performance routing engine with static and parameterized path support
- **Multi-Protocol** — HTTP/1.1, HTTPS/TLS, HTTP/3 (QUIC)
- **WebSocket** — Built-in WebSocket support with automatic protocol upgrade
- **Swagger Auto-Generation** — OpenAPI 3.0 documentation generated from API definitions
- **Response Compression** — Automatic gzip/deflate compression
- **Flexible Caching** — Extensible method-level cache with custom persistence implementations
- **Interceptor Chain** — Ordered HTTP interceptor pipeline with system/user handler separation
