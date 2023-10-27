package log

import "fmt"

var (
	std Logger = &defaultLogger{}
)

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	std.Trace(args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	std.Debug(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	std.Info(args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	std.Warn(args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	std.Error(args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	std.Panic(args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	std.Tracef(format, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

func SetLogger(logger Logger) {
	std = logger
}

func SetLevel(level Level, f func(logger Logger)) {
	f(std)
}

var _ Logger = (*defaultLogger)(nil)

type defaultLogger struct {
}

// Debug implements Logger.
func (*defaultLogger) Debug(args ...interface{}) {
	fmt.Println(args...)
}

// Debugf implements Logger.
func (*defaultLogger) Debugf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// Error implements Logger.
func (*defaultLogger) Error(args ...interface{}) {
	fmt.Println(args...)
}

// Errorf implements Logger.
func (*defaultLogger) Errorf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// Fatal implements Logger.
func (*defaultLogger) Fatal(args ...interface{}) {
	fmt.Println(args...)
}

// Info implements Logger.
func (*defaultLogger) Info(args ...interface{}) {
	fmt.Println(args...)
}

// Infof implements Logger.
func (*defaultLogger) Infof(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// Panic implements Logger.
func (*defaultLogger) Panic(args ...interface{}) {
	fmt.Println(args...)
}

// Trace implements Logger.
func (*defaultLogger) Trace(args ...interface{}) {
	fmt.Println(args...)
}

// Tracef implements Logger.
func (*defaultLogger) Tracef(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

// Warn implements Logger.
func (*defaultLogger) Warn(args ...interface{}) {
	fmt.Println(args...)
}

// Warnf implements Logger.
func (*defaultLogger) Warnf(format string, args ...interface{}) {
	panic("unimplemented")
}
