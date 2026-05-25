# 响应压缩

框架根据客户端的 `Accept-Encoding` 请求头自动压缩 HTTP 响应。

## 支持的压缩算法

| 算法 | Header 值 |
|------|----------|
| gzip | `gzip` |
| deflate | `deflate` |

## 使用方式

无需配置 — 压缩是自动的。客户端通过 `Accept-Encoding` 请求压缩：

```bash
curl -H "Accept-Encoding: gzip" http://127.0.0.1:8080/api/data
```

响应包含：

```
Content-Encoding: gzip
```

## 工作原理

`CompressStd` 拦截器在响应阶段以 `Order = 1500` 运行：

1. 检查 `Accept-Encoding` 是否包含支持的算法
2. 如果响应是 `RetAdapter`（Stream、Html 等）：
   - 创建管道：原始 reader → 压缩器 → pipe writer
   - 通过新的 `Stream` 响应流式输出压缩内容
   - 保留原始状态码和响应头
3. 如果响应是序列化的 JSON（`def.Content`）：
   - 在内存中压缩字节
   - 包装为 `Stream` 响应

## 自定义压缩算法

```go
compress.CompressRegister["br"] = &brotliCompressor{}
```

`Compress` 接口：

```go
type Compress interface {
    New(io.Writer) io.WriteCloser
    ContentEncoding() string
}
```
