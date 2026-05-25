# Error Handling

## Built-in Error Type

The `def.Error` struct provides a standard error response:

```go
type Error struct {
    ErrorMessage string `json:"error"`
    Code         int    `json:"code"`
}
```

```go
// Auto-generated hash code
def.NewError("something went wrong")
// {"error":"something went wrong","code":1234567890}

// Custom error code
def.NewErrorCode("message", 10001)
// {"error":"message","code":10001}
```

## Panic Handling

The framework recovers from panics in handlers and returns a JSON error response:

```go
api.GET(func() interface{} {
    panic("something broke")
}, "/panic")
```

Response:
```json
{
    "error": "something broke",
    "code": 2893576325
}
```

## Custom Error Handlers

Register type-specific error handlers using `api.RegisterErrorHandler`:

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

When a handler panics with a registered error type, the custom handler transforms the error into a response.

## Error Handler Resolution

The framework resolves errors in this order:

1. **Custom handler** — if `reflect.TypeOf(err)` has a registered handler
2. **String** — if the error is a `string`, wraps in `def.NewError(string)`
3. **error interface** — if the error implements `error`, wraps in `def.NewError(err.Error())`
4. **Default** — returns empty string

## Handler Error Response

Error responses are written as:

```go
rw.Header().Add("Content-Type", "application/json;charset=utf-8")
rw.WriteHeader(http.StatusInternalServerError) // 500
rw.Write(jsonBytes)
```

## Not Found (404)

When no route matches, the `NotFind` interceptor returns:

```json
{
    "path": "/unknown/path",
    "msg": "Not find Path"
}
```

With HTTP status `404 Not Found`.
