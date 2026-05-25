# Return Types

Handler functions can return any type. The framework serializes the return value based on its type.

## Default JSON Serialization

Returning a map, struct, or slice produces JSON:

```go
api.GET(func() interface{} {
    return map[string]string{"message": "hello"}
}, "/json")
```

```json
{"message":"hello"}
```

## Plain Text

Returning a `string` sends it as-is without JSON wrapping:

```go
api.GET(func() interface{} {
    return "plain text response"
}, "/text")
```

## No Return Value

Returning `nil` sends an empty response with HTTP 200:

```go
api.GET(func() interface{} {
    // do something
    return nil
}, "/nothing")
```

## Stream Response — File Download

`api.NewStream(io.Reader)` creates a streaming response for file downloads:

```go
api.GET(func() interface{} {
    f, _ := os.Open("large_file.zip")
    return api.NewStream(f).SetName("download.zip")
}, "/download")
```

### Stream Features

| Feature | Method |
|---------|--------|
| Set filename | `.SetName(name)` |
| Rate limiting | `.SetRateLimit(bytesPerSec)` |
| Custom Content-Type | `.SetContentType(type)` |
| Custom status code | `.SetCode(code)` |
| Custom headers | `.AddHeader(key, value)` |

### HTTP Range Requests

Streams automatically support HTTP Range requests when the underlying reader implements `io.Seeker`:

```bash
curl -H "Range: bytes=0-1023" http://127.0.0.1:8080/download
# Returns 206 Partial Content with Content-Range header
```

The stream automatically:
- Detects file size via `io.Seeker`
- Sets `Accept-Ranges: bytes`
- Parses Range header and seeks to the requested position
- Returns `206 Partial Content` with `Content-Range` header

### Rate Limiting

```go
api.GET(func() interface{} {
    f, _ := os.Open("video.mp4")
    return api.NewStream(f).SetRateLimit(1024 * 1024) // 1 MB/s
}, "/download")
```

Rate limiting uses `golang.org/x/time/rate` token bucket algorithm.

## HTML Template Rendering

### Inline Template

```go
api.GET(func() interface{} {
    return api.Html("<h1>Hello {{.Name}}</h1>", map[string]string{"Name": "World"})
}, "/html")
```

### Template from embed.FS

```go
//go:embed views
var views embed.FS

api.GET(func() interface{} {
    data := map[string]interface{}{"Title": "Users", "Users": userList}
    return api.HtmlView(views, "views/list.html", data)
}, "/users")
```

Uses Go's `html/template` for rendering. The template execution has a 1-minute timeout.

## Redirect

```go
api.GET(func() interface{} {
    return api.NewRedirect("https://example.com")
}, "/redirect")
```

Returns `302 Found` with `Location` header.

## Custom Response — api.NewResp

Full control over status code, headers, and body:

```go
api.GET(func() interface{} {
    return api.NewResp(map[string]string{"status": "ok"}).
        SetCode(http.StatusCreated).
        SetHeader(map[string]string{"X-Custom": "value"})
}, "/custom")
```

### Resp Methods

| Method | Description |
|--------|-------------|
| `.SetCode(code)` | Set HTTP status code |
| `.SetHeader(map)` | Set response headers |
| `.SetContentType(type)` | Override Content-Type |
| `.SetSerialize(s)` | Use custom serializer |
| `.SetReader(r)` | Use custom io.Reader as body |

### Convenience Functions

```go
// Status code only
api.Status(http.StatusNoContent)

// Headers only
api.Header(map[string]string{"X-Key": "value"})
```

## Custom Return Handler

Register custom return type handlers using `api.RegisterReturnHandler`:

```go
api.RegisterReturnHandler(&MyRetAdapter{})
```

Interface:

```go
type RetAdapter interface {
    ContentType() string
    Return() io.Reader
    Register() []reflect.Type
}
```

Optional interfaces for richer behavior:

```go
type AppendHeader interface {
    Append(header ReadHeader) map[string]string
}

type HttpStatus interface {
    Code() int
}
```
