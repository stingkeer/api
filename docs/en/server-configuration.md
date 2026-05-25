# Server Configuration

## Configuration Options

Server behavior is configured through optional functional parameters:

```go
api.StartService(
    api.WithListen(":9090"),
)
```

### Available Options

| Option | Type | Description |
|--------|------|-------------|
| `WithListen(addr)` | `string` | Set listen address. Default: `0.0.0.0:8080` |
| `WithTLS(certPEM, keyPEM)` | `[]byte, []byte` | Set TLS certificate from bytes |
| `WithTLSFile(certFile, keyFile)` | `string, string` | Set TLS certificate from file paths |
| `WithTLSConfig(config)` | `*tls.Config` | Set custom TLS configuration |
| `WithCa(caPEM)` | `[]byte` | Set CA certificate |
| `WithDwarfMode(mode)` | `dwarf.FilterMode` | Set DWARF filter mode |
| `WithExePath(path)` | `string` | Set executable file path for DWARF parsing |

## ServerConfig

The `ServerConfig` struct holds all configuration:

```go
type ServerConfig struct {
    // accessed via methods
    listen    string
    dwarf     *dwarf.DwarfMaker
    loadPath  *string
    // TLS fields
    caPEMBlock    []byte
    certPEMBlock  []byte
    keyPEMBlock   []byte
    tlsConfig     *tls.Config
}
```

### Accessors

```go
conf := api.StartService // internal use only

// Read current config (via ServerConfig methods)
config.Listen()          // get listen address
config.Dwarf()           // get DwarfMaker
config.LoadPath()        // get executable path
```

### Modifiers

```go
config.SetListen(":9090")
config.SetDwarfMaker(dwarfMaker)
config.SetDwarfMode(dwarf.RuntimeExclude)
config.SetExePath("/path/to/binary")
config.AddIncludeRegex("^myapp/.+$")
config.AddExclude("github.com")
```

## Start HTTP Server

```go
api.StartService(api.WithListen(":8080"))
```

## Start HTTPS Server

Using certificate bytes:

```go
certPEM, _ := os.ReadFile("cert.pem")
keyPEM, _ := os.ReadFile("key.pem")

api.StartTLSService(
    api.WithListen(":443"),
    api.WithTLS(certPEM, keyPEM),
)
```

Using certificate files:

```go
api.StartTLSService(
    api.WithListen(":443"),
    api.WithTLSFile("cert.pem", "key.pem"),
)
```

Using custom TLS config:

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

## DWARF Configuration

### Include/Exclude Packages

Control which functions are indexed by the DWARF parser:

```go
api.WithDwarfMode(dwarf.SelfInclude) // Only include go.aew.app packages

// Or add custom rules:
config.AddIncludeRegex("^myapp/.+$")  // Include packages matching regex
config.AddExclude("github.com")       // Exclude packages by prefix
```

### Custom Executable Path

By default, DWARF reads the current executable. Override with:

```go
api.StartService(
    api.WithExePath("/path/to/my/app"),
)
```

Or via environment variable:

```bash
API_DLL=/path/to/my/app go run main.go
```

## Test Mode

Test mode is auto-detected and skips DWARF parsing:

- Binary name ends with `.test` AND args contain `-test.` prefix
- Or set `API_TEST=1` to force-enable DWARF even in tests

In test mode, routes register with minimal metadata (no DWARF parameter names).
