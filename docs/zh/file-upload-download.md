# 文件上传与下载

## 文件上传

使用 `multipart.Reader` 作为函数参数处理文件上传：

```go
api.POST(func(reader multipart.Reader) interface{} {
    form, _ := reader.ReadForm(0)
    fHeader := form.File["file"][0]
    f, _ := fHeader.Open()
    bytes, _ := io.ReadAll(f)
    return bytes
}, "/upload")
```

```bash
curl -X POST http://127.0.0.1:8080/upload \
  -F 'file=@/path/to/file.txt'
```

### 多文件上传

```go
api.POST(func(reader multipart.Reader) interface{} {
    form, _ := reader.ReadForm(32 << 20) // 最大 32MB 内存
    var files []string
    for _, fh := range form.File["files"] {
        files = append(files, fh.Filename)
    }
    return files
}, "/upload-multi")
```

### 读取表单字段

```go
api.POST(func(reader multipart.Reader) interface{} {
    form, _ := reader.ReadForm(0)
    return map[string]string{
        "name":  form.Value["name"][0],
        "email": form.Value["email"][0],
    }
}, "/form")
```

## 文件下载

使用 `api.NewStream` 配合 `os.File` 进行文件下载：

```go
api.GET(func() interface{} {
    f, _ := os.Open("report.pdf")
    return api.NewStream(f).SetName("report.pdf")
}, "/download")
```

### 自动功能

当底层 reader 为 `*os.File`（实现了 `io.Seeker`）时：

| 功能 | 自动行为 |
|------|---------|
| Content-Type | 从前 512 字节自动检测 |
| Content-Length | 从文件大小设置 |
| Accept-Ranges | 设置为 `bytes` |
| Content-Disposition | `attachment; filename=report.pdf` |
| Range 支持 | 返回 `206 Partial Content` |

### HTTP Range 请求

```bash
curl -H "Range: bytes=0-1023" http://127.0.0.1:8080/download
# 206 Partial Content
# Content-Range: bytes 0-1023/146515
# Content-Length: 1024
```

### 限速下载

```go
api.GET(func() interface{} {
    f, _ := os.Open("video.mp4")
    return api.NewStream(f).
        SetName("video.mp4").
        SetRateLimit(1024 * 1024) // 1 MB/s
}, "/download")
```

速率限制使用令牌桶算法（`golang.org/x/time/rate`）。

### 自定义流式响应

```go
api.GET(func() interface{} {
    data := generateLargeData()
    reader := bytes.NewReader(data)
    return api.NewStream(reader).
        SetContentType("application/pdf").
        SetName("report.pdf")
}, "/report")
```

## 流式代理（非文件）

任何 `io.Reader` 都可以流式传输：

```go
api.GET(func() interface{} {
    resp, _ := http.Get("https://example.com/large-file")
    return api.NewStream(resp.Body)
}, "/proxy")
```
