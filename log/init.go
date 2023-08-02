package log

func NewGlobal(cc *Config) error {
	// init global core, the global core's level must be DEBUG.
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
