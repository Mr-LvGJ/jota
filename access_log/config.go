package access_log

import (
	"github.com/Mr-LvGJ/jota/log"
)

var accessLog = &AccessLogger{}

type Option func(*AccessLogger)

type AccessLogger struct {
	loggerConfig *log.Config
	*log.Logger
}

var defaultLogConfig = &log.Config{
	Level:      "info",
	CallerSkip: 3,

	Structured: true,
	CoreLevel:  "debug",
	ErrorFile:  false,

	Filename:   "",
	MaxSize:    10,
	MaxAge:     14,
	MaxBackups: 14,
	LocalTime:  true,
	Compress:   true,
}

func WithLogConfig(cfg *log.Config) Option {
	return func(c *AccessLogger) {
		c.loggerConfig = cfg
	}
}

func WithLogger(logger *log.Logger) Option {
	return func(c *AccessLogger) {
		c.Logger = logger
	}
}

func NewConfig(opts ...Option) *AccessLogger {
	for _, opt := range opts {
		opt(accessLog)
	}
	if &accessLog.loggerConfig == nil {
		accessLog.loggerConfig = defaultLogConfig
	}
	if accessLog.Logger != nil {
		return accessLog
	}
	logger, err := log.NewLogger(accessLog.loggerConfig)
	if err != nil {
		panic(err)
	}
	accessLog.Logger = logger
	return accessLog
}
