# File Upload & Download

## File Upload

Use `multipart.Reader` as a function parameter to handle file uploads:

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

### Multiple Files

```go
api.POST(func(reader multipart.Reader) interface{} {
    form, _ := reader.ReadForm(32 << 20) // 32MB max memory
    var files []string
    for _, fh := range form.File["files"] {
        files = append(files, fh.Filename)
    }
    return files
}, "/upload-multi")
```

### Read Form Fields

```go
api.POST(func(reader multipart.Reader) interface{} {
    form, _ := reader.ReadForm(0)
    return map[string]string{
        "name":  form.Value["name"][0],
        "email": form.Value["email"][0],
    }
}, "/form")
```

## File Download

Use `api.NewStream` with an `os.File` for file downloads:

```go
api.GET(func() interface{} {
    f, _ := os.Open("report.pdf")
    return api.NewStream(f).SetName("report.pdf")
}, "/download")
```

### Automatic Features

When the underlying reader is an `*os.File` (implements `io.Seeker`):

| Feature | Automatic Behavior |
|---------|-------------------|
| Content-Type | Auto-detected from first 512 bytes |
| Content-Length | Set from file size |
| Accept-Ranges | `bytes` |
| Content-Disposition | `attachment; filename=report.pdf` |
| Range support | Returns `206 Partial Content` |

### HTTP Range Request

```bash
curl -H "Range: bytes=0-1023" http://127.0.0.1:8080/download
# 206 Partial Content
# Content-Range: bytes 0-1023/146515
# Content-Length: 1024
```

### Rate-Limited Download

```go
api.GET(func() interface{} {
    f, _ := os.Open("video.mp4")
    return api.NewStream(f).
        SetName("video.mp4").
        SetRateLimit(1024 * 1024) // 1 MB/s
}, "/download")
```

Rate limiting uses a token bucket algorithm (`golang.org/x/time/rate`).

### Custom Stream

```go
api.GET(func() interface{} {
    data := generateLargeData()
    reader := bytes.NewReader(data)
    return api.NewStream(reader).
        SetContentType("application/pdf").
        SetName("report.pdf")
}, "/report")
```

## Streaming Response (Non-File)

Any `io.Reader` can be streamed:

```go
api.GET(func() interface{} {
    resp, _ := http.Get("https://example.com/large-file")
    return api.NewStream(resp.Body)
}, "/proxy")
```
