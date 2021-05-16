package logger

// An Option configures a Logger
type Option interface {
	Apply(Logger) Logger
}

// OptionFunc is a function that configures a Logger
type OptionFunc func(Logger) Logger

// Apply is a function that set value to Logger
func (f OptionFunc) Apply(engine Logger) Logger {
	return f(engine)
}

func WithFields(fields map[string]interface{}) Option {
	return OptionFunc(func(engine Logger) Logger {
		zlog := engine.With().Fields(fields).Logger()
		return Logger{zlog}
	})
}

func WithStr(key, value string) Option {
	return OptionFunc(func(engine Logger) Logger {
		zlog := engine.With().Str(key, value).Logger()
		return Logger{zlog}
	})
}
