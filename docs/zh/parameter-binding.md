# 参数绑定

框架自动将请求数据绑定到函数参数。参数名通过编译后二进制中嵌入的 **DWARF 调试信息** 获取，无需结构体标签或手动注解。

## 工作原理

1. 启动时，`DwarfMaker` 读取可执行文件的 DWARF 调试段
2. 对每个注册的 handler，`LookFun(fn)` 返回函数的参数名和类型
3. 请求到达时，`callerDefault.Call()` 使用存储的元信息将请求数据映射到函数参数

## 基本类型

所有 Go 基本类型会自动从查询字符串值转换：

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

支持的类型：`bool`、`string`、`int`、`int8`、`int16`、`int32`、`int64`、`uint`、`uint8`、`uint16`、`uint32`、`uint64`、`float32`、`float64`

缺失的参数会使用零值（空字符串、0、false）。

## 必填参数

使用 `def.XxxReq` 类型声明必填参数。如果请求中缺少该参数，框架会 panic 返回错误。

```go
api.GET(func(name string, password def.StringReq) interface{} {
    return "ok"
}, "/login")
```

可用必填类型：
- `def.StringReq` — 必填字符串
- `def.IntReq` — 必填 int
- `def.Int8Req` — 必填 int8
- `def.Int16Req` — 必填 int16
- `def.Int32Req` — 必填 int32
- `def.Int64Req` — 必填 int64

## 结构体参数

当函数参数为结构体类型时，框架通过 `json` tag 从查询参数绑定字段：

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

支持嵌套结构体 — 框架递归绑定结构体字段。

## Body 参数

当参数的 DWARF 名称为 `"body"` 时，框架从请求体中读取数据并反序列化：

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

Body 参数规则：
- **每个 handler 只能有一个 body 参数** — 超过 1 个会 panic
- 支持：结构体（JSON 解码）、字符串（原始内容）、`[]byte`（原始字节）、其他类型的切片

## 文件上传 (multipart.Reader)

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

## *http.Request — 原始请求

```go
api.GET(func(req *http.Request) interface{} {
    return req.URL.String()
}, "/raw")
```

也支持 `http.Request`（非指针）— 接收请求结构体的副本。

## def.Header — 请求头

```go
api.GET(func(h def.Header) interface{} {
    return h.Values("Accept-Encoding")
}, "/headers")
```

`def.Header` 接口提供：

```go
type Header interface {
    Cookie                           // SetCookie, Cookie
    ReadHeader                       // Get(key), Values(key)
    WriteHeader                      // Add(key, value)
}
```

写操作（`Add`、`SetCookie`）写入的是**响应**头，不是请求头。

## big.Int — 大整数

```go
api.GET(func(amount big.Int) interface{} {
    return amount.String()
}, "/big")
```

支持 `big.Int`（值）和 `*big.Int`（指针）。

## WebSocket 参数

详见 [WebSocket](websocket.md) 中的 `ws.WSCtx` 参数类型。

## 自定义类型适配器

使用 `api.RegisterTypeMapper` 注册自定义参数类型适配器：

```go
api.RegisterTypeMapper(&MyAdapter{})
```

适配器接口：

```go
type Adapter interface {
    Mapper(param *ParamWarp) reflect.Value
    Register() []reflect.Type
}
```

- `Register()` 返回此适配器处理的类型列表
- `Mapper()` 将请求数据转换为 `reflect.Value`

对于泛型类型，使用 `call.RegisterGenericTypeMapper`。
