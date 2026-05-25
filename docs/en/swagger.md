# Swagger Documentation

The framework can automatically generate OpenAPI 3.0 documentation from your API definitions.

## Building with Swagger

```bash
go build -tags 'swagger'
```

Access the Swagger UI at: `http://localhost:8080/ui/index.html#/`

## Route Metadata

Add metadata to individual routes using `.Swagger()`:

```go
api.GET(func(name string) interface{} {
    return "hello " + name
}, "/greet").Swagger(func(sw def.SwaggerOps) {
    sw.SetSummary("Greeting endpoint")
    sw.SetTag("user")
    sw.SetDescription("Returns a personalized greeting")
    sw.SetParameterDescription("name", "User's display name")
})
```

### SwaggerOps Methods

| Method | Description | OpenAPI Field |
|--------|-------------|---------------|
| `SetSummary(text)` | Short summary | `summary` |
| `SetDescription(text)` | Detailed description | `description` |
| `SetTag(tag)` | Group tag | `tags` |
| `SetParameterDescription(name, desc)` | Parameter description | `parameter.description` |

## Security Schemes

Define security requirements for route groups:

```go
api.Routes(
    api.GET(handler, "/api/data"),
).Swagger(func(sw def.SwaggerSecurity) {
    sw.SecuritJwt("bearerAuth")              // JWT Bearer
    sw.SecuritCookie("auth", "session_id")   // Cookie auth
    sw.SecuritApiHeader("key", "X-API-Key")  // API Key header
})
```

### Generated Security Schemes

**JWT:**
```json
{
  "bearerAuth": {
    "type": "http",
    "scheme": "bearer",
    "bearerformat": "JWT"
  }
}
```

**Cookie:**
```json
{
  "auth": {
    "type": "apiKey",
    "in": "cookie",
    "name": "session_id"
  }
}
```

**API Key Header:**
```json
{
  "key": {
    "type": "apiKey",
    "in": "header",
    "name": "X-API-Key"
  }
}
```

## Automatic Schema Generation

The Swagger generator automatically:

1. **Parameter detection** — infers parameter location (query, path, body) from the function signature
2. **Type mapping** — maps Go types to OpenAPI types:
   - `string` → `string`
   - `int`, `int64` → `integer` + `int64`
   - `int32` → `integer` + `int32`
   - `float64` → `number`
   - `bool` → `boolean`
   - `struct` → `object` with `$ref`
   - `[]T` → `array` with items
3. **Required detection** — `def.XxxReq` types and path parameters are marked required
4. **Schema references** — struct types are extracted into `components/schemas`
5. **Path parameter conversion** — `<name>` in URLs becomes `{name}` in OpenAPI paths

## Generated Output

The generator produces an OpenAPI 3.0.3 spec:

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

## Environment Variable

Set the server URL via environment variable:

```bash
api.listen=http://myserver.com go build -tags 'swagger'
```
