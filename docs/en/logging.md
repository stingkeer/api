# Logging

The framework provides a replaceable logging interface.

## Default Logger

The default logger writes to `stdout` via `fmt.Println`:

```go
log.Info("server started")
log.Infof("listening on %s", addr)
log.Warn("slow request")
log.Warnf("param %s not found", name)
log.Error(err)
log.Errorf("connection failed: %v", err)
```

## Log Levels

From highest to lowest severity:

| Level | Description |
|-------|-------------|
| `PanicLevel` | Logs then calls panic |
| `FatalLevel` | Logs then exits with status 1 |
| `ErrorLevel` | Errors that should be noted |
| `WarnLevel` | Non-critical warnings |
| `InfoLevel` | General operational entries |
| `DebugLevel` | Verbose debugging info |
| `TraceLevel` | Finer-grained than Debug |

## Replace Logger

Use your own logger implementation:

```go
import "go.aew.app/api.v1/log"

type MyLogger struct{}

func (l *MyLogger) Info(args ...interface{})            { /* ... */ }
func (l *MyLogger) Infof(format string, args ...interface{}) { /* ... */ }
func (l *MyLogger) Warn(args ...interface{})            { /* ... */ }
func (l *MyLogger) Warnf(format string, args ...interface{}) { /* ... */ }
func (l *MyLogger) Error(args ...interface{})           { /* ... */ }
func (l *MyLogger) Errorf(format string, args ...interface{}) { /* ... */ }
func (l *MyLogger) Debug(args ...interface{})           { /* ... */ }
func (l *MyLogger) Debugf(format string, args ...interface{}) { /* ... */ }
func (l *MyLogger) Trace(args ...interface{})           { /* ... */ }
func (l *MyLogger) Tracef(format string, args ...interface{}) { /* ... */ }
func (l *MyLogger) Panic(args ...interface{})           { /* ... */ }
func (l *MyLogger) Fatal(args ...interface{})           { /* ... */ }

func init() {
    log.SetLogger(&MyLogger{})
}
```

## Framework Log Messages

The framework logs at these points:

| Context | Level | Message |
|---------|-------|---------|
| Route registration | Info | `[GET] pkg.func(arg1,arg2) mapping url = /path` |
| Server start | Info | `listen addr :8080` |
| DWARF init | Print | `DwarfMaker init use Xms len N` |
| Request trace | Trace | `incoming req HttpMethod [GET], Url [/path]` |
| No match | Trace | `not match /path` |
| Method mismatch | Warn | `not support HttpMethod DELETE` |
| Unsupported param | Warn | `not support type X` |
| DWARF lookup fail | Error | `not find name [pkg.func]` |
