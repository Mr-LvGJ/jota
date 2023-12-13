package access_log

import (
	"github.com/Mr-LvGJ/jota/log"
)

var accessLog = &AccessLog{}

type option func(*AccessLog)

type AccessLog struct {
	loggerConfig *log.Config
	logger       *log.Logger
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

func WithLogConfig(cfg *log.Config) option {
	return func(c *AccessLog) {
		c.loggerConfig = cfg
	}
}

func WithLogger(logger *log.Logger) option {
	return func(c *AccessLog) {
		c.logger = logger
	}
}

func newConfig(opts ...option) {
	for _, opt := range opts {
		opt(accessLog)
	}
	if &accessLog.loggerConfig == nil {
		accessLog.loggerConfig = defaultLogConfig
	}
	if accessLog.logger != nil {
		return
	}
	logger, err := log.NewLogger(accessLog.loggerConfig)
	if err != nil {
		panic(err)
	}
	accessLog.logger = logger
}
