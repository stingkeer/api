# 快速开始

## 安装

```bash
go get go.aew.app/api.v1
```

## Hello World

```go
package main

import (
    "go.aew.app/api.v1"
)

func main() {
    api.GET(func() interface{} {
        return "hello world"
    }, "/hello")

    api.StartService(api.WithListen(":8080"))
}
```

运行：

```bash
go run main.go
# 输出: [GET] main.main.func1() mapping url = /hello
#       listen addr :8080
```

测试：

```bash
curl http://127.0.0.1:8080/hello
# 响应: hello world
```

## 从 Query String 绑定参数

框架自动将 URL 查询参数绑定到函数参数。参数名通过 DWARF 调试信息获取，无需手动解析或添加结构体标签。

```go
api.GET(func(username, password string) interface{} {
    return map[string]string{
        "username": username,
        "password": password,
    }
}, "/login")
```

```bash
curl 'http://127.0.0.1:8080/login?username=alice&password=123456'
# 响应: {"password":"123456","username":"alice"}
```

## 多个路由

```go
func main() {
    api.GET(func() interface{} {
        return "hello"
    }, "/")

    api.GET(func(name string) interface{} {
        return "hello " + name
    }, "/greet")

    api.POST(func(body LoginRequest) interface{} {
        return body
    }, "/login")

    api.StartService(api.WithListen(":8080"))
}
```

## HTTPS 服务

```go
api.StartTLSService(
    api.WithListen(":443"),
    api.WithTLSFile("cert.pem", "key.pem"),
)
```

## 构建 Swagger UI

```bash
go build -tags 'swagger'
```

访问 `http://localhost:8080/ui/index.html#/` 查看文档。

## 构建 HTTP/3 (QUIC)

```bash
go build -tags 'http3'
```
