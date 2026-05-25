# 日志系统

框架提供可替换的日志接口。

## 默认日志器

默认日志器通过 `fmt.Println` 输出到 stdout：

```go
log.Info("server started")
log.Infof("listening on %s", addr)
log.Warn("slow request")
log.Warnf("param %s not found", name)
log.Error(err)
log.Errorf("connection failed: %v", err)
```

## 日志级别

从高到低严重程度：

| 级别 | 说明 |
|------|------|
| `PanicLevel` | 记录后调用 panic |
| `FatalLevel` | 记录后以状态码 1 退出 |
| `ErrorLevel` | 应被注意的错误 |
| `WarnLevel` | 非关键警告 |
| `InfoLevel` | 常规操作信息 |
| `DebugLevel` | 详细调试信息 |
| `TraceLevel` | 比 Debug 更细粒度 |

## 替换日志器

使用自定义日志器实现：

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

## 框架日志输出点

| 上下文 | 级别 | 消息 |
|--------|------|------|
| 路由注册 | Info | `[GET] pkg.func(arg1,arg2) mapping url = /path` |
| 服务启动 | Info | `listen addr :8080` |
| DWARF 初始化 | Print | `DwarfMaker init use Xms len N` |
| 请求追踪 | Trace | `incoming req HttpMethod [GET], Url [/path]` |
| 无匹配路由 | Trace | `not match /path` |
| 方法不匹配 | Warn | `not support HttpMethod DELETE` |
| 不支持的参数类型 | Warn | `not support type X` |
| DWARF 查找失败 | Error | `not find name [pkg.func]` |
