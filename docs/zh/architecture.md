# 架构设计

## 请求处理流水线

```
HTTP Request
    │
    ├─── Order=0 ─── 系统级前置处理器（CORS 等）
    │
    ├─── Order 1-99 ── 系统处理器
    │   ├── Static (99) — 提供嵌入文件
    │   └── (可扩展)
    │
    ├─── Order 100-999 ── 框架核心
    │   ├── API Handler (100)
    │   │   ├── 路由匹配（基数树）
    │   │   ├── 参数绑定（DWARF + 类型适配器）
    │   │   ├── 中间件执行
    │   │   ├── Handler 调用（带方法代理链）
    │   │   └── 响应序列化
    │   └── API Response — 写入序列化数据
    │
    ├─── 响应阶段 ──
    │   ├── Compress (1500) — gzip/deflate 压缩
    │   └── NotFind (MaxUint) — 404 JSON 响应
    │
    └─── HTTP Response
```

## 核心组件

### ServerMain

入口点，实现 `http.Handler`。委托给 `http.DoHttp()`。

### Context (def.DefaultContext)

中央配置，持有：
- `Match` — 基数树路由器
- `Pool` — 方法元信息池
- `Caller` — 参数绑定 + 调用（默认 WsCaller）
- `Serialize` — JSON 序列化器

### DWARF 解析器（dwarf 包）

读取编译后二进制的调试信息，提取函数参数名。

```
编译 → 带有 DWARF 的二进制
    │
    ├── DwarfMaker.Init()
    │   ├── macOS: macho.Open() → DWARF()
    │   ├── Linux: elf.Open() → DWARF()
    │   └── Windows: pe.Open() → DWARF()
    │
    ├── 遍历 DWARF entries
    │   ├── TagSubprogram → 函数名
    │   └── TagFormalParameter → 参数名
    │
    └── 过滤运行时/标准库包
```

### 基数树路由器（match 包）

O(k) 路径匹配（k = 路径长度）。支持：
- 静态段 — 精确字节匹配
- 参数段 (`<name>`) — 捕获非 `/` 字符
- 正则段 (`<name:pattern>`) — 匹配编译后的正则表达式
- 通配符段 (`<name:.*>`) — 匹配剩余路径

### 参数绑定（call 包）

```
Caller.Call(entry, request)
    │
    ├── 从 Pool 获取 MethodInfo
    │
    ├── 遍历每个参数：
    │   ├── name == "body" → 读取请求体，反序列化
    │   ├── 类型有注册 Adapter → 使用注册的适配器
    │   ├── 类型有泛型 Adapter → 使用泛型适配器
    │   ├── 类型为结构体 → 通过 json tag 从 Query 绑定
    │   └── 未知类型 → 设置零值
    │
    ├── 执行中间件链
    │
    └── 通过方法代理链调用 → handler 函数
```

### 方法代理链

```
SetMethodProxy(proxy1)
SetMethodProxy(proxy2)

调用顺序：
proxy2 → proxy1 → reflect.Value.Call()

使用者：
- 缓存模块（最外层代理）
- 用户定义的 AOP 代理
```

### 类型适配器（call/types）

| 适配器 | 类型 | 数据来源 |
|--------|------|---------|
| `BaseType` | bool, string, int*, uint*, float* | 查询字符串 |
| `TypeRequire` | def.IntReq, def.StringReq 等 | 查询字符串（必填） |
| `TypeRequireG` | def.Int[any], def.String[any] 等 | 查询字符串（必填） |
| `HttpType` | *http.Request, http.Request | 注入 |
| `HeadType` | def.Header | 注入 |
| `WSType` | ws.WSCtx, *ws.WSCtx | WebSocket 升级 |
| `FileType` | multipart.Reader | 请求体 |
| `BigType` | big.Int, *big.Int | 查询字符串 |

### 返回适配器（call/rettypes）

| 适配器 | Content-Type | 功能 |
|--------|-------------|------|
| `Stream` | 自动检测 | Range 请求、速率限制、文件下载 |
| `Html` | text/html | Go 模板渲染、embed.FS 支持 |
| `Redirect` | text/html | 302 重定向 |
| `Resp` | application/json | 自定义状态码、响应头、序列化器 |

## 包结构

```
api.v1/
├── api.go          → 公开 API（GET、POST 等）
├── def.go          → 公开类型（NewStream、Html 等）
├── config.go       → ServerConfig + 配置函数
├── http.go         → StartService（HTTP/1.1）
├── http3.go        → StartTLSService（HTTP/3，构建标签）
├── middleware.go    → Routes()、AddRoutes()、Middleware()
├── websocket.go    → WebSocket 使用示例
├── respose.go      → Status()、Header()
│
├── def/            → 核心接口和类型
│   ├── http.go     → HttpMethod、MiddleWare、Option、Header、ContentType
│   ├── open_inf.go → Request、Serialize、Match、Caller
│   ├── open_entry.go → MethodInfo、Entry、ParamWarp、MethodsPools
│   ├── call_adapter.go → Adapter、RetAdapter
│   ├── require.go  → Int[T]、String[T] 泛型必填类型
│   ├── require_base.go → IntReq、StringReq 类型别名
│   ├── error.go    → Error 结构体
│   ├── swagger.go  → SwaggerSecurity、SwaggerOps
│   ├── order.go    → HandlerOrder 常量
│   └── flush.go    → Flusher 接口
│
├── http/           → HTTP 层
│   ├── do.go       → DoHttp()、拦截器管道
│   ├── api.go      → ApiIntercept、ApiResponse 处理器
│   ├── errors.go   → 错误处理器注册
│   ├── ret_adapter.go → RetAdapter 注册
│   ├── static.go   → 静态文件服务
│   ├── notfind.go  → 404 处理器
│   └── header.go   → readHeader 辅助
│
├── call/           → 参数绑定和调用
│   ├── init.go     → 类型适配器注册
│   ├── caller_default.go → 默认参数绑定逻辑
│   ├── trace_caller.go → 带中间件的调用器
│   ├── ws_celler.go → WebSocket 感知的调用器
│   ├── default_invoke.go → 方法代理链
│   ├── call_adapter.go → RealCall 接口
│   ├── types/      → 参数类型适配器
│   └── rettypes/   → 返回类型适配器
│
├── kit/            → 引导和扩展
│   ├── core/       → HttpM()、option、swagger、context 初始化
│   ├── swagger/    → OpenAPI 3.0 生成
│   ├── ws/         → WebSocket 上下文和连接池
│   └── handler/    → Compress、CORS 处理器
│
├── dwarf/          → DWARF 调试信息解析器
├── match/          → 基数树路由器
├── cache/          → 方法级缓存
├── intercept/      → HTTP 拦截器接口
├── serialize/      → JSON 序列化
├── log/            → 可替换日志器
└── utils/          → 泛型 sync.Map、工具函数
```
