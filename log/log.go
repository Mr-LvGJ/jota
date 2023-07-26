package log

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const contextLogKey = iota

var gLogger *Logger

func Global() *Logger {
	return gLogger
}

func SetGlobal(l *Logger) {
	if gLogger != nil {
		if err := gLogger.Sync(); err != nil && !strings.Contains(err.Error(), "sync /dev/stdout") {
			l.Warn(context.Background(), "close global logger failed", "err", err)
		}
	}
	gLogger = l
}

func NewLogger(config *Config, cores ...zapcore.Core) (*Logger, error) {
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, err
	}
	return &Logger{
		zl: zap.New(
			zapcore.NewTee(cores...),
			zap.WithCaller(true),
			zap.AddCallerSkip(config.CallerSkip),
			zap.AddStacktrace(zapcore.ErrorLevel),
		),
		l: level,
	}, nil
}

func Debug(ctx context.Context, msg string, fields ...interface{}) {
	gLogger.Debug(ctx, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...interface{}) {
	gLogger.Info(ctx, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...interface{}) {
	gLogger.Warn(ctx, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...interface{}) {
	gLogger.Error(ctx, msg, fields...)
}

func Panic(ctx context.Context, msg string, fields ...interface{}) {
	gLogger.Fatal(ctx, msg, fields...)
}

func WithValue(ctx context.Context, key string, value interface{}) context.Context {
	return gLogger.WithValues(ctx, key, value)
}

type Logger struct {
	zl *zap.Logger
	l  zapcore.Level
}

func (l *Logger) Sync() (err error) {
	return l.zl.Sync()
}

func (l *Logger) Enabled(level zapcore.Level) bool {
	return l.l.Enabled(level)
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	l.Log(ctx, zapcore.DebugLevel, msg, fields...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...interface{}) {
	l.Log(ctx, zapcore.InfoLevel, msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	l.Log(ctx, zapcore.WarnLevel, msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...interface{}) {
	l.Log(ctx, zapcore.ErrorLevel, msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...interface{}) {
	l.Log(ctx, zapcore.FatalLevel, msg, fields...)
}

func (l *Logger) Panic(ctx context.Context, msg string, fields ...interface{}) {
	l.Log(ctx, zapcore.PanicLevel, msg, fields...)
}

func (l *Logger) Log(ctx context.Context, level zapcore.Level, msg string, kvs ...interface{}) {
	if !l.Enabled(level) {
		return
	}

	ctxFields, _ := ctx.Value(contextLogKey).([]zap.Field)
	fields := append(ctxFields, zapFields(kvs)...)

	switch level {
	case zapcore.DebugLevel:
		l.zl.Debug(msg, fields...)
	case zapcore.InfoLevel:
		l.zl.Info(msg, fields...)
	case zapcore.WarnLevel:
		l.zl.Warn(msg, fields...)
	case zapcore.ErrorLevel:
		l.zl.Error(msg, fields...)
	case zapcore.FatalLevel:
		l.zl.Fatal(msg, fields...)
	case zapcore.PanicLevel:
		l.zl.Fatal(msg, fields...)
	}
}

func (l *Logger) WithValues(ctx context.Context, kvs ...interface{}) context.Context {
	v, _ := ctx.Value(contextLogKey).([]zap.Field)
	return context.WithValue(ctx, contextLogKey, append(v, zapFields(kvs)...))
}
