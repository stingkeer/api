# 中间件

中间件函数在请求到达 handler 之前拦截 HTTP 请求。通过返回非 nil 值可以短路请求。

## 路由级中间件

使用 `api.Routes()` 或 `api.AddRoutes()` 为一组路由应用中间件：

```go
authMiddleware := func(req *http.Request) (ret any) {
    token := req.Header.Get("Authorization")
    if token == "" {
        return api.NewResp("unauthorized").SetCode(http.StatusUnauthorized)
    }
    return nil // nil = 继续执行下一个中间件或 handler
}

api.Routes(
    api.GET(func() interface{} {
        return "protected data"
    }, "/data1"),
    api.GET(func() interface{} {
        return "more data"
    }, "/data2"),
).Middleware(authMiddleware)
```

## 中间件行为

- **返回 `nil`** — 继续执行下一个中间件或 handler
- **返回非 nil** — 短路调用链，返回值作为响应发送

中间件按添加顺序执行：

```go
api.Routes(
    api.GET(handler, "/api/data"),
).Middleware(
    loggingMiddleware,   // 第一个执行
    authMiddleware,      // 第二个执行
    rateLimitMiddleware, // 第三个执行
)
```

## Swagger 安全定义

为路由组关联安全方案，用于 Swagger 文档：

```go
api.Routes(
    api.GET(handler, "/api/data"),
).Swagger(func(sw def.SwaggerSecurity) {
    sw.SecuritJwt("bearerAuth")
})
```

可用安全类型：

| 方法 | 说明 | Swagger 方案 |
|------|------|-------------|
| `SecuritJwt(name)` | JWT Bearer 令牌 | `http` + `bearer` |
| `SecuritCookie(name, cookieName)` | Cookie 认证 | `apiKey` + `cookie` |
| `SecuritApiHeader(tag, headerName)` | Header API Key | `apiKey` + `header` |

## HTTP 拦截器

更低级别的 HTTP 拦截，实现 `intercept.HttpIntercept` 接口：

```go
type HttpIntercept interface {
    Http(rw http.ResponseWriter, req *http.Request, ctx *HttpContext) bool
    Order() def.HandlerOrder
}
```

- 返回 `true` 停止调用链
- 返回 `false` 继续执行
- `Order()` 决定执行优先级

使用 `api.AddHttpHandle()` 注册，Order 值必须 >= 100。

### Handler Order 范围

| Order 范围 | 用途 |
|-----------|------|
| 0 | 系统级前置处理器（如 CORS） |
| 1–99 | 系统处理器（如静态文件 = 99） |
| 100–999 | 框架核心（如 API 路由 = 100） |
| 1000–1499 | 用户后置处理器 |
| 1500+ | 后处理（如压缩 = 1500） |
| MaxUint | 404 兜底 |

## 方法代理（AOP）

方法级拦截使用 `call.SetMethodProxy`：

```go
call.SetMethodProxy(func(fn call.MethodCaller, m *def.MethodInfo, args []reflect.Value) []reflect.Value {
    // 前置通知
    result := fn.Invoke(m, args)
    // 后置通知
    return result
})
```

方法代理形成链式调用 — 可以叠加多个代理。缓存模块通过此机制实现透明缓存。
