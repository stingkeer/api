# 错误处理

## 内建错误类型

`def.Error` 结构体提供标准错误响应：

```go
type Error struct {
    ErrorMessage string `json:"error"`
    Code         int    `json:"code"`
}
```

```go
// 自动生成哈希错误码
def.NewError("something went wrong")
// {"error":"something went wrong","code":1234567890}

// 自定义错误码
def.NewErrorCode("message", 10001)
// {"error":"message","code":10001}
```

## Panic 处理

框架从 handler 的 panic 中恢复并返回 JSON 错误响应：

```go
api.GET(func() interface{} {
    panic("something broke")
}, "/panic")
```

响应：
```json
{
    "error": "something broke",
    "code": 2893576325
}
```

## 自定义错误处理器

使用 `api.RegisterErrorHandler` 注册类型特定的错误处理器：

```go
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

api.RegisterErrorHandler(reflect.TypeOf(&ValidationError{}), func(err interface{}) interface{} {
    e := err.(*ValidationError)
    return api.NewResp(e).SetCode(http.StatusBadRequest)
})
```

当 handler 以注册的错误类型 panic 时，自定义处理器将错误转换为响应。

## 错误处理器解析顺序

框架按以下顺序解析错误：

1. **自定义处理器** — 如果 `reflect.TypeOf(err)` 有注册的处理器
2. **字符串** — 如果错误是 `string`，包装为 `def.NewError(string)`
3. **error 接口** — 如果错误实现了 `error`，包装为 `def.NewError(err.Error())`
4. **默认** — 返回空字符串

## 错误响应格式

错误响应写入为：

```go
rw.Header().Add("Content-Type", "application/json;charset=utf-8")
rw.WriteHeader(http.StatusInternalServerError) // 500
rw.Write(jsonBytes)
```

## 404 未找到

当没有路由匹配时，`NotFind` 拦截器返回：

```json
{
    "path": "/unknown/path",
    "msg": "Not find Path"
}
```

HTTP 状态码 `404 Not Found`。
