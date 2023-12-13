package log

type Config struct {
	Level      string `default:"debug"`
	CallerSkip int    `default:"3"`

	// ZapCore config
	Structured bool   `default:"true"`
	CoreLevel  string `default:"debug"'`
	ErrorFile  bool   `default:"false"`
	Prod       bool

	Filename   string
	MaxSize    int `default:"500"`
	MaxAge     int `default:"7"`
	MaxBackups int `default:"5"`
	LocalTime  bool
	Compress   bool `default:"true"`
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
