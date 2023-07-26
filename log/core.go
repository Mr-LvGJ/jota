package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var globalCore zapcore.Core

func GlobalCore() zapcore.Core {
	return globalCore
}

func SetGlobalCore(c zapcore.Core) {
	globalCore = c
}

func NewZapCore(cfg *Config) (zapcore.Core, error) {
	return newCore(cfg)
}

func newCore(cfg *Config) (zapcore.Core, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(cfg.CoreLevel)); err != nil {
		return nil, err
	}
	return zapcore.NewCore(newEncoder(cfg), newWriter(cfg), level), nil
}

func newEncoder(cfg *Config) zapcore.Encoder {
	var encoderCfg zapcore.EncoderConfig
	if cfg.Prod {
		encoderCfg = zap.NewProductionEncoderConfig()
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}
	if cfg.Structured {
		return zapcore.NewJSONEncoder(encoderCfg)
	} else {
		return zapcore.NewConsoleEncoder(encoderCfg)
	}
}

func newWriter(cfg *Config) zapcore.WriteSyncer {
	if len(cfg.Filename) > 0 {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
	} else {
		return zapcore.AddSync(os.Stdout)
	}
	return nil
}
