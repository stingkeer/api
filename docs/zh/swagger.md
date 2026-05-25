# Swagger 文档

框架可以根据 API 定义自动生成 OpenAPI 3.0 文档。

## 构建 Swagger 版本

```bash
go build -tags 'swagger'
```

访问 Swagger UI：`http://localhost:8080/ui/index.html#/`

## 路由元信息

通过 `.Swagger()` 为单个路由添加文档元信息：

```go
api.GET(func(name string) interface{} {
    return "hello " + name
}, "/greet").Swagger(func(sw def.SwaggerOps) {
    sw.SetSummary("问候接口")
    sw.SetTag("user")
    sw.SetDescription("根据名字返回个性化问候")
    sw.SetParameterDescription("name", "用户显示名称")
})
```

### SwaggerOps 方法

| 方法 | 说明 | OpenAPI 字段 |
|------|------|-------------|
| `SetSummary(text)` | 简短摘要 | `summary` |
| `SetDescription(text)` | 详细描述 | `description` |
| `SetTag(tag)` | 分组标签 | `tags` |
| `SetParameterDescription(name, desc)` | 参数描述 | `parameter.description` |

## 安全方案

为路由组定义安全要求：

```go
api.Routes(
    api.GET(handler, "/api/data"),
).Swagger(func(sw def.SwaggerSecurity) {
    sw.SecuritJwt("bearerAuth")              // JWT Bearer
    sw.SecuritCookie("auth", "session_id")   // Cookie 认证
    sw.SecuritApiHeader("key", "X-API-Key")  // API Key Header
})
```

### 生成的安全方案

**JWT：**
```json
{
  "bearerAuth": {
    "type": "http",
    "scheme": "bearer",
    "bearerformat": "JWT"
  }
}
```

**Cookie：**
```json
{
  "auth": {
    "type": "apiKey",
    "in": "cookie",
    "name": "session_id"
  }
}
```

**API Key Header：**
```json
{
  "key": {
    "type": "apiKey",
    "in": "header",
    "name": "X-API-Key"
  }
}
```

## 自动 Schema 生成

Swagger 生成器自动完成：

1. **参数检测** — 从函数签名推断参数位置（query、path、body）
2. **类型映射** — Go 类型到 OpenAPI 类型：
   - `string` → `string`
   - `int`、`int64` → `integer` + `int64`
   - `int32` → `integer` + `int32`
   - `float64` → `number`
   - `bool` → `boolean`
   - `struct` → `object` + `$ref`
   - `[]T` → `array` + items
3. **必填检测** — `def.XxxReq` 类型和路径参数标记为 required
4. **Schema 引用** — 结构体类型提取到 `components/schemas`
5. **路径参数转换** — URL 中的 `<name>` 变为 OpenAPI 路径中的 `{name}`

## 生成输出

生成器产出 OpenAPI 3.0.3 规范：

```json
{
  "openapi": "3.0.3",
  "info": {
    "title": "Golang API Generate",
    "description": "This is a sample server for api"
  },
  "servers": [{"url": "http://localhost:8080"}],
  "paths": { ... },
  "components": {
    "schemas": { ... },
    "securitySchemes": { ... }
  }
}
```

## 环境变量

通过环境变量设置服务器 URL：

```bash
api.listen=http://myserver.com go build -tags 'swagger'
```
