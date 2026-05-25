# api.v1 Documentation

A declarative Go web framework that infers parameter binding from function signatures using DWARF debug information.

- [English Documentation](en/)
- [中文文档](zh/)

## Overview

| Feature | Description |
|---------|-------------|
| Declarative Routing | Function signatures define parameter binding automatically |
| DWARF Parameter Resolution | Reads compiled binary DWARF info for zero-reflection parameter names |
| Radix Tree Router | High-performance routing with static and parameterized paths |
| Multi-Protocol | HTTP/1.1, HTTPS/TLS, HTTP/3 (QUIC) |
| WebSocket | Built-in WebSocket with automatic protocol upgrade |
| Swagger | Auto-generated OpenAPI 3.0 documentation |
| Compression | Automatic gzip/deflate response compression |
| Caching | Extensible method-level cache with custom persistence |
| Interceptor Pipeline | Ordered HTTP interceptor chain with system/user separation |
