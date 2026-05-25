# 路由注册

框架提供标准 HTTP 方法注册函数：

```go
api.HEAD(handler, "/path")
api.GET(handler, "/path")
api.POST(handler, "/path")
api.PUT(handler, "/path")
api.PATCH(handler, "/path")
api.DELETE(handler, "/path")
api.OPTIONS(handler, "/path")
```

## 静态路径

```go
api.GET(func() interface{} {
    return "users list"
}, "/users")
```

## 参数化路径

使用 `<name>` 语法定义路径参数，参数名必须与函数参数名一致。

```go
api.GET(func(id string) interface{} {
    return map[string]string{"id": id}
}, "/user/<id>")
```

```bash
curl http://127.0.0.1:8080/user/42
# 响应: {"id":"42"}
```

多个参数：

```go
api.GET(func(userId, postId string) interface{} {
    return map[string]string{
        "user": userId,
        "post": postId,
    }
}, "/user/<userId>/post/<postId>")
```

## 正则约束路径

使用 `<name:pattern>` 语法添加正则约束：

```go
api.GET(func(id string) interface{} {
    return map[string]string{"id": id}
}, "/user/<id:[0-9]+>")
```

此路由仅匹配数字 ID，请求 `/user/abc` 不会匹配。

## 通配符路径

使用 `<name:.*>` 匹配任意子路径：

```go
api.GET(func(path string) interface{} {
    return map[string]string{"path": path}
}, "/files/<path:.*>")
```

```bash
curl http://127.0.0.1:8080/files/a/b/c.txt
# 响应: {"path":"a/b/c.txt"}
```

## 路由组与中间件

使用 `api.Routes()` 或 `api.AddRoutes()` 将路由分组并应用中间件：

```go
api.Routes(
    api.GET(func() interface{} { return "data1" }, "/api/data1"),
    api.GET(func() interface{} { return "data2" }, "/api/data2"),
).Middleware(authMiddleware)
```

详见 [中间件](middleware.md)。

## Swagger 元信息

通过 `.Swagger()` 为单个路由添加文档元信息：

```go
api.GET(func(name string) interface{} {
    return "hello " + name
}, "/greet").Swagger(func(sw def.SwaggerOps) {
    sw.SetSummary("问候接口")
    sw.SetTag("user")
    sw.SetDescription("根据名字返回问候语")
    sw.SetParameterDescription("name", "用户名称")
})
```

## 路由匹配原理

框架使用**基数树** (`match/radix_tree.go`) 实现 O(k) 路径匹配（k 为路径长度）。支持：

- 静态段 — 精确字节匹配
- 参数段 (`<name>`) — 捕获非 `/` 字符
- 正则段 (`<name:pattern>`) — 匹配编译后的正则表达式
- 通配符段 (`<name:.*>`) — 匹配剩余路径

路由注册在初始化时通过 `HttpM()` 完成：
1. 将路由添加到基数树
2. 读取函数的 DWARF 调试信息获取参数名
3. 将 `MethodInfo` 存入 `MethodsPools` 供运行时使用
