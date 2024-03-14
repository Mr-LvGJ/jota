package log

func init() {
	if err := NewGlobal(&Config{Level: "info", CallerSkip: 3}); err != nil {
		panic(err)
	}
}

func NewGlobal(cc *Config) error {
	// init global core, the global core's level must be DEBUG.
	c := DefaultConfig()
	if cc.Level == "" {
		cc.Level = c.Level
	}
	if cc.CallerSkip == 0 {
		cc.CallerSkip = c.CallerSkip
	}
	if cc.CoreLevel == "" {
		cc.CoreLevel = c.CoreLevel
	}
	core, err := NewZapCore(cc)
	if err != nil {
		return err
	}
	SetGlobalCore(core)

	logger, err := NewLoggerWithCore(cc, core)
	if err != nil {
		return err
	}
	SetGlobal(logger)
	return nil
}
