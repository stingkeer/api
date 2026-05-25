# Routing

The framework provides functions for all standard HTTP methods:

```go
api.HEAD(handler, "/path")
api.GET(handler, "/path")
api.POST(handler, "/path")
api.PUT(handler, "/path")
api.PATCH(handler, "/path")
api.DELETE(handler, "/path")
api.OPTIONS(handler, "/path")
```

## Static Paths

```go
api.GET(func() interface{} {
    return "users list"
}, "/users")
```

## Parameterized Paths

Use `<name>` syntax to define path parameters. The parameter name must match a function argument name.

```go
api.GET(func(id string) interface{} {
    return map[string]string{"id": id}
}, "/user/<id>")
```

```bash
curl http://127.0.0.1:8080/user/42
# Response: {"id":"42"}
```

Multiple parameters:

```go
api.GET(func(userId, postId string) interface{} {
    return map[string]string{
        "user": userId,
        "post": postId,
    }
}, "/user/<userId>/post/<postId>")
```

## Regex-Constrained Paths

Use `<name:pattern>` syntax to add regex constraints:

```go
api.GET(func(id string) interface{} {
    return map[string]string{"id": id}
}, "/user/<id:[0-9]+>")
```

This route only matches numeric IDs. A request to `/user/abc` will not match.

## Wildcard Paths

Use `<name:.*>` to match any sub-path:

```go
api.GET(func(path string) interface{} {
    return map[string]string{"path": path}
}, "/files/<path:.*>")
```

```bash
curl http://127.0.0.1:8080/files/a/b/c.txt
# Response: {"path":"a/b/c.txt"}
```

## Route Groups with Middleware

Use `api.Routes()` or `api.AddRoutes()` to group routes and apply middleware:

```go
api.Routes(
    api.GET(func() interface{} { return "data1" }, "/api/data1"),
    api.GET(func() interface{} { return "data2" }, "/api/data2"),
).Middleware(authMiddleware)
```

See [Middleware](middleware.md) for details.

## Swagger Metadata

Attach Swagger metadata to individual routes using `.Swagger()`:

```go
api.GET(func(name string) interface{} {
    return "hello " + name
}, "/greet").Swagger(func(sw def.SwaggerOps) {
    sw.SetSummary("Greeting endpoint")
    sw.SetTag("user")
    sw.SetDescription("Returns a greeting message")
    sw.SetParameterDescription("name", "User name")
})
```

## How Routing Works

The framework uses a **radix tree** (`match/radix_tree.go`) for O(k) path matching where k is the path length. The tree supports:

- Static segments — exact byte matching
- Parameter segments (`<name>`) — capture non-`/` characters
- Regex segments (`<name:pattern>`) — match against compiled regex
- Wildcard segments (`<name:.*>`) — match remaining path

Route registration happens at init time via `HttpM()`, which:
1. Adds the route to the radix tree
2. Reads the function's DWARF debug info for parameter names
3. Stores the `MethodInfo` in `MethodsPools` for runtime lookup
