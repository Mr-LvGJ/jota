package log

type Config struct {
	Level      string
	CallerSkip int

	// ZapCore config
	Structured bool
	CoreLevel  string
	ErrorFile  bool
	Prod       bool
	
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
}

var defaultConfig = &Config{
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

func DefaultConfig() *Config {
	return defaultConfig
}
