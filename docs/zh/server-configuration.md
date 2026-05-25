# 服务启动配置

## 配置选项

通过可选函数参数配置服务行为：

```go
api.StartService(
    api.WithListen(":9090"),
)
```

### 可用配置项

| 配置函数 | 类型 | 说明 |
|---------|------|------|
| `WithListen(addr)` | `string` | 设置监听地址，默认 `0.0.0.0:8080` |
| `WithTLS(certPEM, keyPEM)` | `[]byte, []byte` | 从字节设置 TLS 证书 |
| `WithTLSFile(certFile, keyFile)` | `string, string` | 从文件路径设置 TLS 证书 |
| `WithTLSConfig(config)` | `*tls.Config` | 设置自定义 TLS 配置 |
| `WithCa(caPEM)` | `[]byte` | 设置 CA 证书 |
| `WithDwarfMode(mode)` | `dwarf.FilterMode` | 设置 DWARF 过滤模式 |
| `WithExePath(path)` | `string` | 设置 DWARF 解析的可执行文件路径 |

## ServerConfig 结构

`ServerConfig` 持有所有配置：

```go
type ServerConfig struct {
    listen        string            // 监听地址
    dwarf         *dwarf.DwarfMaker // DWARF 解析器
    loadPath      *string           // 可执行文件路径
    caPEMBlock    []byte            // CA 证书
    certPEMBlock  []byte            // TLS 证书
    keyPEMBlock   []byte            // TLS 私钥
    tlsConfig     *tls.Config       // 自定义 TLS 配置
}
```

### 访问器

```go
config.Listen()    // 获取监听地址
config.Dwarf()     // 获取 DwarfMaker
config.LoadPath()  // 获取可执行文件路径
```

### 修改器

```go
config.SetListen(":9090")
config.SetDwarfMaker(dwarfMaker)
config.SetDwarfMode(dwarf.RuntimeExclude)
config.SetExePath("/path/to/binary")
config.AddIncludeRegex("^myapp/.+$")
config.AddExclude("github.com")
```

## 启动 HTTP 服务

```go
api.StartService(api.WithListen(":8080"))
```

## 启动 HTTPS 服务

使用证书字节：

```go
certPEM, _ := os.ReadFile("cert.pem")
keyPEM, _ := os.ReadFile("key.pem")

api.StartTLSService(
    api.WithListen(":443"),
    api.WithTLS(certPEM, keyPEM),
)
```

使用证书文件：

```go
api.StartTLSService(
    api.WithListen(":443"),
    api.WithTLSFile("cert.pem", "key.pem"),
)
```

使用自定义 TLS 配置：

```go
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
    // ...
}

api.StartTLSService(
    api.WithListen(":443"),
    api.WithTLSConfig(tlsConfig),
)
```

## DWARF 配置

### 包含/排除包

控制 DWARF 解析器索引哪些函数：

```go
api.WithDwarfMode(dwarf.SelfInclude) // 仅包含 go.aew.app 包

// 或添加自定义规则：
config.AddIncludeRegex("^myapp/.+$")  // 包含匹配正则的包
config.AddExclude("github.com")       // 按前缀排除包
```

### 自定义可执行文件路径

默认 DWARF 读取当前可执行文件。可通过配置覆盖：

```go
api.StartService(
    api.WithExePath("/path/to/my/app"),
)
```

或通过环境变量：

```bash
API_DLL=/path/to/my/app go run main.go
```

## 测试模式

测试模式自动检测并跳过 DWARF 解析：

- 可执行文件名以 `.test` 结尾且参数包含 `-test.` 前缀
- 或设置 `API_TEST=1` 强制在测试中启用 DWARF

测试模式下，路由使用最少的元数据注册（无 DWARF 参数名）。
