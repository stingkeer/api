# Architecture

## Request Processing Pipeline

```
HTTP Request
    │
    ├─── Order=0 ─── System Pre-processors (CORS, etc.)
    │
    ├─── Order 1-99 ── System Handlers
    │   ├── Static (99) — serve embedded files
    │   └── (extensible)
    │
    ├─── Order 100-999 ── Framework Core
    │   ├── API Handler (100)
    │   │   ├── Route matching (radix tree)
    │   │   ├── Parameter binding (DWARF + type adapters)
    │   │   ├── Middleware execution
    │   │   ├── Handler invocation (with method proxy chain)
    │   │   └── Response serialization
    │   └── API Response — write serialized data
    │
    ├─── Response Phase ──
    │   ├── Compress (1500) — gzip/deflate
    │   └── NotFind (MaxUint) — 404 JSON response
    │
    └─── HTTP Response
```

## Core Components

### ServerMain

Entry point implementing `http.Handler`. Delegates to `http.DoHttp()`.

### Context (def.DefaultContext)

Central configuration holding:
- `Match` — radix tree router
- `Pool` — method metadata pool
- `Caller` — parameter binding + invocation (WsCaller by default)
- `Serialize` — JSON serializer

### DWARF Parser (dwarf package)

Reads compiled binary's debug information to extract function parameter names.

```
Compile → Binary with DWARF
    │
    ├── DwarfMaker.Init()
    │   ├── macOS: macho.Open() → DWARF()
    │   ├── Linux: elf.Open() → DWARF()
    │   └── Windows: pe.Open() → DWARF()
    │
    ├── Iterate DWARF entries
    │   ├── TagSubprogram → function names
    │   └── TagFormalParameter → parameter names
    │
    └── Filter out runtime/stdlib packages
```

### Radix Tree Router (match package)

O(k) path matching where k = path length. Supports:
- Static segments — exact byte matching
- Parameter segments (`<name>`) — capture non-`/` characters
- Regex segments (`<name:pattern>`) — match against compiled regex
- Wildcard segments (`<name:.*>`) — match remaining path

### Parameter Binding (call package)

```
Caller.Call(entry, request)
    │
    ├── Get MethodInfo from pool
    │
    ├── For each parameter:
    │   ├── name == "body" → read request body, deserialize
    │   ├── type has Adapter → use registered adapter
    │   ├── type has generic Adapter → use generic adapter
    │   ├── type is struct → bind from query via json tags
    │   └── unknown → set zero value
    │
    ├── Execute middleware chain
    │
    └── Invoke via method proxy chain → handler function
```

### Method Proxy Chain

```
SetMethodProxy(proxy1)
SetMethodProxy(proxy2)

Invocation order:
proxy2 → proxy1 → reflect.Value.Call()

Used by:
- Cache module (outermost proxy)
- User-defined AOP proxies
```

### Type Adapters (call/types)

| Adapter | Types | Source |
|---------|-------|--------|
| `BaseType` | bool, string, int*, uint*, float* | Query string |
| `TypeRequire` | def.IntReq, def.StringReq, etc. | Query string (required) |
| `TypeRequireG` | def.Int[any], def.String[any], etc. | Query string (required) |
| `HttpType` | *http.Request, http.Request | Injected |
| `HeadType` | def.Header | Injected |
| `WSType` | ws.WSCtx, *ws.WSCtx | WebSocket upgrade |
| `FileType` | multipart.Reader | Request body |
| `BigType` | big.Int, *big.Int | Query string |

### Return Adapters (call/rettypes)

| Adapter | Content-Type | Features |
|---------|-------------|----------|
| `Stream` | auto-detected | Range requests, rate limiting, file download |
| `Html` | text/html | Go template rendering, embed.FS support |
| `Redirect` | text/html | 302 redirect |
| `Resp` | application/json | Custom status, headers, serializer |

## Package Map

```
api.v1/
├── api.go          → Public API (GET, POST, etc.)
├── def.go          → Public types (NewStream, Html, etc.)
├── config.go       → ServerConfig + option functions
├── http.go         → StartService (HTTP/1.1)
├── http3.go        → StartTLSService (HTTP/3, build tag)
├── middleware.go    → Routes(), AddRoutes(), Middleware()
├── websocket.go    → WebSocket usage examples
├── respose.go      → Status(), Header()
│
├── def/            → Core interfaces and types
│   ├── http.go     → HttpMethod, MiddleWare, Option, Header, ContentType
│   ├── open_inf.go → Request, Serialize, Match, Caller
│   ├── open_entry.go → MethodInfo, Entry, ParamWarp, MethodsPools
│   ├── call_adapter.go → Adapter, RetAdapter
│   ├── require.go  → Int[T], String[T] generic required types
│   ├── require_base.go → IntReq, StringReq type aliases
│   ├── error.go    → Error struct
│   ├── swagger.go  → SwaggerSecurity, SwaggerOps
│   ├── order.go    → HandlerOrder constants
│   └── flush.go    → Flusher interface
│
├── http/           → HTTP layer
│   ├── do.go       → DoHttp(), interceptor pipeline
│   ├── api.go      → ApiIntercept, ApiResponse handlers
│   ├── errors.go   → Error handler registration
│   ├── ret_adapter.go → RetAdapter registration
│   ├── static.go   → Static file serving
│   ├── notfind.go  → 404 handler
│   └── header.go   → readHeader helper
│
├── call/           → Parameter binding and invocation
│   ├── init.go     → Type adapter registration
│   ├── caller_default.go → Default parameter binding logic
│   ├── trace_caller.go → Middleware-aware caller
│   ├── ws_celler.go → WebSocket-aware caller
│   ├── default_invoke.go → Method proxy chain
│   ├── call_adapter.go → RealCall interface
│   ├── types/      → Parameter type adapters
│   └── rettypes/   → Return type adapters
│
├── kit/            → Bootstrapping and extensions
│   ├── core/       → HttpM(), option, swagger, context init
│   ├── swagger/    → OpenAPI 3.0 generation
│   ├── ws/         → WebSocket context and pool
│   └── handler/    → Compress, CORS handlers
│
├── dwarf/          → DWARF debug info parser
├── match/          → Radix tree router
├── cache/          → Method-level caching
├── intercept/      → HTTP interceptor interface
├── serialize/      → JSON serialization
├── log/            → Replaceable logger
└── utils/          → Generic sync.Map, helpers
```
