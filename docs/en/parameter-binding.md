# Parameter Binding

The framework automatically binds request data to function arguments. Parameter names are obtained from **DWARF debug information** embedded in the compiled binary — no struct tags or manual annotation needed for parameter name resolution.

## How It Works

1. At startup, `DwarfMaker` reads the executable's DWARF debug section
2. For each registered handler, `LookFun(fn)` returns the function's parameter names and types
3. At request time, `callerDefault.Call()` maps request data to function arguments using the stored metadata

## Basic Types

All Go basic types are automatically converted from query string values:

```go
api.GET(func(name string, age int, score float64, active bool) interface{} {
    return map[string]interface{}{
        "name":   name,
        "age":    age,
        "score":  score,
        "active": active,
    }
}, "/info")
```

```bash
curl 'http://127.0.0.1:8080/info?name=Alice&age=30&score=95.5&active=true'
```

Supported types: `bool`, `string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`

Missing parameters receive zero values (empty string, 0, false).

## Required Parameters

Use `def.XxxReq` types to declare required parameters. If the parameter is missing from the request, the framework panics with an error message.

```go
api.GET(func(name string, password def.StringReq) interface{} {
    return "ok"
}, "/login")
```

Available required types:
- `def.StringReq` — required string
- `def.IntReq` — required int
- `def.Int8Req` — required int8
- `def.Int16Req` — required int16
- `def.Int32Req` — required int32
- `def.Int64Req` — required int64

## Struct Parameters

When a function argument is a struct type, the framework binds fields from query parameters using `json` tags:

```go
type SearchQuery struct {
    Keyword string `json:"keyword"`
    Page    int    `json:"page"`
    Size    int    `json:"size"`
}

api.GET(func(q SearchQuery) interface{} {
    return q
}, "/search")
```

```bash
curl 'http://127.0.0.1:8080/search?keyword=golang&page=1&size=10'
```

Nested structs are supported — the framework recursively binds struct fields.

## Body Parameters

When a parameter's DWARF name is `"body"`, the framework reads the request body and deserializes it:

```go
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

api.POST(func(body LoginRequest) interface{} {
    return body
}, "/login")
```

```bash
curl -X POST http://127.0.0.1:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123"}'
```

Body parameters:
- **Only one body parameter** is allowed per handler — having 2+ causes a panic
- Supports: struct (JSON decode), string (raw body), `[]byte` (raw bytes), slice of other types

## File Upload (multipart.Reader)

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

## *http.Request — Raw Request Access

```go
api.GET(func(req *http.Request) interface{} {
    return req.URL.String()
}, "/raw")
```

Also supports `http.Request` (non-pointer) — receives a copy of the request struct.

## def.Header — Request Headers

```go
api.GET(func(h def.Header) interface{} {
    return h.Values("Accept-Encoding")
}, "/headers")
```

`def.Header` interface provides:

```go
type Header interface {
    Cookie                           // SetCookie, Cookie
    ReadHeader                       // Get(key), Values(key)
    WriteHeader                      // Add(key, value)
}
```

Write operations (`Add`, `SetCookie`) write to the **response** headers, not the request.

## big.Int — Large Integers

```go
api.GET(func(amount big.Int) interface{} {
    return amount.String()
}, "/big")
```

Supports both `big.Int` (value) and `*big.Int` (pointer).

## WebSocket Parameter

See [WebSocket](websocket.md) for the `ws.WSCtx` parameter type.

## Custom Type Adapters

Register custom parameter type adapters using `api.RegisterTypeMapper`:

```go
api.RegisterTypeMapper(&MyAdapter{})
```

Adapter interface:

```go
type Adapter interface {
    Mapper(param *ParamWarp) reflect.Value
    Register() []reflect.Type
}
```

- `Register()` returns the types this adapter handles
- `Mapper()` converts request data to a `reflect.Value`

For generic types, use `call.RegisterGenericTypeMapper` instead.
