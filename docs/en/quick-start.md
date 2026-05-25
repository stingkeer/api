# Quick Start

## Installation

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

Run:

```bash
go run main.go
# Output: [GET] main.main.func1() mapping url = /hello
#         listen addr :8080
```

Test:

```bash
curl http://127.0.0.1:8080/hello
# Response: hello world
```

## Parameters from Query String

The framework automatically binds URL query parameters to function arguments. Parameter names are obtained from DWARF debug information — no struct tags or manual parsing needed.

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
# Response: {"password":"123456","username":"alice"}
```

## Multiple Routes

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

## HTTPS Server

```go
api.StartTLSService(
    api.WithListen(":443"),
    api.WithTLSFile("cert.pem", "key.pem"),
)
```

## Build with Swagger UI

```bash
go build -tags 'swagger'
```

Access at `http://localhost:8080/ui/index.html#/`

## Build with HTTP/3 (QUIC)

```bash
go build -tags 'http3'
```
