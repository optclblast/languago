package logger

import (
	"languago/infrastructure/logger/wrappers"
)

type (
	abstractLoggerConfig interface {
		GetLogger() Logger
	}

	FormattedLogger interface {
		Writef(format string, args ...any)
		Tracef(format string, args ...any)
		Warnf(format string, args ...any)
		Errorf(format string, args ...any)
		Debugf(format string, args ...any)
		Infof(format string, args ...any)
		Logf(format string, args ...any)
		Panicf(format string, args ...any)
	}

	Logger interface {
		FormattedLogger

		// Level() Level
		// SetLevel(level Level)
		// Backend() *Backend

		Write(args ...any)
		Trace(args ...any)
		Warn(args ...any)
		Error(args ...any)
		Debug(args ...any)
		Info(args ...any)
		Log(args ...any)
		Panic(args ...any)
	}

	logger struct {
		lvl       Level // atomic
		tag       string
		b         *Backend
		writeChan chan<- logEntry
	}

	logEntry struct {
		log   []byte
		level Level
	}

	LogFields map[string]interface{}
)

func ProvideLogger(cfg abstractLoggerConfig) Logger {
	if cfg == nil {
		return provideDefaultLogger(false)
	}

	logger := cfg.GetLogger()
	if logger == nil {
		return provideDefaultLogger(false)
	}
	return logger
}

func provideDefaultLogger(dbg bool) Logger {
	return wrappers.NewZerologWrapper(dbg, wrappers.EnvParam_LOCAL)
}

func LogFieldPair(key string, val any) LogFields {
	lf := make(LogFields, 0)
	lf[key] = val
	return lf
}
