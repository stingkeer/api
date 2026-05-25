# Middleware

Middleware functions intercept HTTP requests before they reach the handler. They can short-circuit the request by returning a non-nil value.

## Route-Level Middleware

Use `api.Routes()` or `api.AddRoutes()` to apply middleware to a group of routes:

```go
authMiddleware := func(req *http.Request) (ret any) {
    token := req.Header.Get("Authorization")
    if token == "" {
        return api.NewResp("unauthorized").SetCode(http.StatusUnauthorized)
    }
    return nil // nil = continue to next middleware/handler
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

## Middleware Behavior

- **Return `nil`** — continue to the next middleware or handler
- **Return non-nil** — short-circuit the chain, the return value is sent as the response

Middleware executes in the order they are added:

```go
api.Routes(
    api.GET(handler, "/api/data"),
).Middleware(
    loggingMiddleware,   // runs first
    authMiddleware,      // runs second
    rateLimitMiddleware, // runs third
)
```

## Swagger Security Definitions

Associate security schemes with routes for Swagger documentation:

```go
api.Routes(
    api.GET(handler, "/api/data"),
).Swagger(func(sw def.SwaggerSecurity) {
    sw.SecuritJwt("bearerAuth")
})
```

Available security types:

| Method | Description | Swagger Scheme |
|--------|-------------|---------------|
| `SecuritJwt(name)` | JWT Bearer token | `http` + `bearer` |
| `SecuritCookie(name, cookieName)` | Cookie-based auth | `apiKey` + `cookie` |
| `SecuritApiHeader(tag, headerName)` | API key in header | `apiKey` + `header` |

## HTTP Interceptor

For lower-level HTTP interception, implement `intercept.HttpIntercept`:

```go
type HttpIntercept interface {
    Http(rw http.ResponseWriter, req *http.Request, ctx *HttpContext) bool
    Order() def.HandlerOrder
}
```

- Return `true` to stop the chain
- Return `false` to continue
- `Order()` determines execution priority

Register with `api.AddHttpHandle()`. The Order value must be >= 100.

### Handler Order Ranges

| Order Range | Usage |
|-------------|-------|
| 0 | System pre-processors (e.g., CORS) |
| 1–99 | System handlers (e.g., Static files = 99) |
| 100–999 | Framework core (e.g., API routing = 100) |
| 1000–1499 | User post-processors |
| 1500+ | Post-processing (e.g., Compression = 1500) |
| MaxUint | 404 fallback |

## Method Proxy (AOP)

For method-level interception (AOP-style), use `call.SetMethodProxy`:

```go
call.SetMethodProxy(func(fn call.MethodCaller, m *def.MethodInfo, args []reflect.Value) []reflect.Value {
    // before advice
    result := fn.Invoke(m, args)
    // after advice
    return result
})
```

Method proxies form a chain — multiple proxies can be stacked. The cache module uses this mechanism.
