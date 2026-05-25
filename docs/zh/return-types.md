# 返回类型

Handler 函数可以返回任意类型。框架根据返回值类型进行序列化。

## 默认 JSON 序列化

返回 map、struct 或 slice 时生成 JSON：

```go
api.GET(func() interface{} {
    return map[string]string{"message": "hello"}
}, "/json")
```

```json
{"message":"hello"}
```

## 纯文本

返回 `string` 时直接发送，不进行 JSON 包装：

```go
api.GET(func() interface{} {
    return "plain text response"
}, "/text")
```

## 无返回值

返回 `nil` 发送空响应，HTTP 200：

```go
api.GET(func() interface{} {
    // do something
    return nil
}, "/nothing")
```

## 流式响应 — 文件下载

`api.NewStream(io.Reader)` 创建流式响应用于文件下载：

```go
api.GET(func() interface{} {
    f, _ := os.Open("large_file.zip")
    return api.NewStream(f).SetName("download.zip")
}, "/download")
```

### Stream 功能

| 功能 | 方法 |
|------|------|
| 设置文件名 | `.SetName(name)` |
| 速率限制 | `.SetRateLimit(bytesPerSec)` |
| 自定义 Content-Type | `.SetContentType(type)` |
| 自定义状态码 | `.SetCode(code)` |
| 自定义响应头 | `.AddHeader(key, value)` |

### HTTP Range 请求

当底层 reader 实现了 `io.Seeker` 时，Stream 自动支持 HTTP Range 请求：

```bash
curl -H "Range: bytes=0-1023" http://127.0.0.1:8080/download
# 返回 206 Partial Content，带 Content-Range 头
```

自动行为：
- 通过 `io.Seeker` 检测文件大小
- 设置 `Accept-Ranges: bytes`
- 解析 Range 头并 seek 到请求位置
- 返回 `206 Partial Content` 和 `Content-Range` 头

### 速率限制

```go
api.GET(func() interface{} {
    f, _ := os.Open("video.mp4")
    return api.NewStream(f).SetRateLimit(1024 * 1024) // 1 MB/s
}, "/download")
```

速率限制使用 `golang.org/x/time/rate` 令牌桶算法。

## HTML 模板渲染

### 内联模板

```go
api.GET(func() interface{} {
    return api.Html("<h1>Hello {{.Name}}</h1>", map[string]string{"Name": "World"})
}, "/html")
```

### 使用 embed.FS 渲染模板

```go
//go:embed views
var views embed.FS

api.GET(func() interface{} {
    data := map[string]interface{}{"Title": "Users", "Users": userList}
    return api.HtmlView(views, "views/list.html", data)
}, "/users")
```

使用 Go 的 `html/template` 渲染，模板执行超时时间为 1 分钟。

## 重定向

```go
api.GET(func() interface{} {
    return api.NewRedirect("https://example.com")
}, "/redirect")
```

返回 `302 Found`，带 `Location` 响应头。

## 自定义响应 — api.NewResp

完全控制状态码、响应头和内容：

```go
api.GET(func() interface{} {
    return api.NewResp(map[string]string{"status": "ok"}).
        SetCode(http.StatusCreated).
        SetHeader(map[string]string{"X-Custom": "value"})
}, "/custom")
```

### Resp 方法

| 方法 | 说明 |
|------|------|
| `.SetCode(code)` | 设置 HTTP 状态码 |
| `.SetHeader(map)` | 设置响应头 |
| `.SetContentType(type)` | 覆盖 Content-Type |
| `.SetSerialize(s)` | 使用自定义序列化器 |
| `.SetReader(r)` | 使用自定义 io.Reader 作为响应体 |

### 快捷函数

```go
// 仅设置状态码
api.Status(http.StatusNoContent)

// 仅设置响应头
api.Header(map[string]string{"X-Key": "value"})
```

## 自定义返回处理器

使用 `api.RegisterReturnHandler` 注册自定义返回类型处理器：

```go
api.RegisterReturnHandler(&MyRetAdapter{})
```

接口：

```go
type RetAdapter interface {
    ContentType() string
    Return() io.Reader
    Register() []reflect.Type
}
```

可选接口实现更丰富的行为：

```go
type AppendHeader interface {
    Append(header ReadHeader) map[string]string
}

type HttpStatus interface {
    Code() int
}
```
