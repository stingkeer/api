# api.v1 - Go 声明式 Web 框架

一个基于 Go 的声明式 Web API 框架，通过函数签名自动推断参数绑定，使用 DWARF 调试信息实现零反射性能损耗的方法参数名解析。

## 目录

1. [快速开始](quick-start.md)
2. [路由注册](routing.md)
3. [参数绑定](parameter-binding.md)
4. [返回类型](return-types.md)
5. [中间件](middleware.md)
6. [静态文件服务](static-files.md)
7. [WebSocket](websocket.md)
8. [文件上传与下载](file-upload-download.md)
9. [缓存系统](caching.md)
10. [服务启动配置](server-configuration.md)
11. [HTTP/3 (QUIC)](http3.md)
12. [Swagger 文档](swagger.md)
13. [错误处理](error-handling.md)
14. [响应压缩](compression.md)
15. [日志系统](logging.md)
16. [架构设计](architecture.md)

## 核心特性

- **声明式路由** — 函数签名自动定义 URL 参数绑定
- **DWARF 参数解析** — 读取编译后二进制的 DWARF 调试信息获取参数名，避免运行时反射开销
- **基数树路由** — 高性能路由匹配引擎，支持静态路径和参数化路径
- **多协议支持** — HTTP/1.1、HTTPS/TLS、HTTP/3 (QUIC)
- **WebSocket** — 内建 WebSocket 支持，自动协议升级
- **Swagger 自动生成** — 根据 API 定义自动生成 OpenAPI 3.0 文档
- **响应压缩** — 自动 gzip/deflate 压缩
- **灵活缓存** — 可扩展的方法级缓存，支持自定义持久化实现
- **拦截器链** — 可排序的 HTTP 拦截器管道，系统与用户处理器分离
