package logger

import (
	"log"
	"log/slog"

	"github.com/sirupsen/logrus"
)

type (
	abstractLoggerConfig interface {
		GetLogger() Logger
	}

	DefaultLogger struct {
		dbgMode bool
		log     *log.Logger
	}

	LogrusWrapper struct {
		dbgMode bool
		log     *logrus.Entry
	}

	SlogLogger struct {
		dbgMode bool
		log     *slog.Logger
	}

	Logger interface {
		Warn(kv ...interface{})
		Debug(kv ...interface{})
		Info(kv ...interface{})
		Log(kv ...interface{})
	}
)

func NewLogrusWrapper(dbg bool) *LogrusWrapper {
	return &LogrusWrapper{
		dbgMode: dbg,
		log:     logrus.NewEntry(logrus.New()),
	}
}

func NewSLogLogger(dbg bool) *SlogLogger {
	return &SlogLogger{
		dbgMode: dbg,
		log:     slog.Default(),
	}
}

func NewDefaultLogger(dbg bool) *DefaultLogger {
	return provideDefaultLogger(dbg)
}

// TODO
// Extends std Logger type to fit the Logger interface
func (l *DefaultLogger) Warn(kv ...interface{}) {
	l.log.SetPrefix("[WARN]")
	l.log.Println(kv...)
}
func (l *DefaultLogger) Err(kv ...interface{}) {
	l.log.SetPrefix("[ERROR]")
	l.log.Println(kv...)
}
func (l *DefaultLogger) Debug(kv ...interface{}) {
	if !l.dbgMode {
		return
	}
	l.log.SetPrefix("[DEBUG]")
	l.log.Println(kv...)
}
func (l *DefaultLogger) Info(kv ...interface{}) {
	l.log.SetPrefix("[INFO]")
	l.log.Println(kv...)
}
func (l *DefaultLogger) Log(kv ...interface{}) {
	l.log.SetPrefix("[LOG]")
	l.log.Println(kv...)
}
func (l *DefaultLogger) Panic(kv ...interface{}) {
	l.log.SetPrefix("[PANIC]")
	l.log.Panicln(kv...)
}

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

func WithPrefix(log Logger, kv ...string) Logger {
	entry := log.(*LogrusWrapper).log

	if len(kv) == 1 {
		entry = entry.WithField("module", kv[0])
	} else if len(kv)%2 == 0 {
		for i := 0; i < len(kv); i += 2 {
			entry = entry.WithField(kv[i], kv[i+1])
		}
	}

	return &LogrusWrapper{
		log: entry,
	}
}

func provideDefaultLogger(dbg bool) *DefaultLogger {
	logger := DefaultLogger{
		dbgMode: dbg,
		log:     log.Default(),
	}
	logger.log.SetFlags(log.LstdFlags)
	logger.log.Println("Error getting logger config. Default logger set.")
	return &logger
}

func (l *LogrusWrapper) Warn(kv ...interface{}) {
	l.log.Warnln(kv...)
}

func (l *LogrusWrapper) Debug(kv ...interface{}) {
	if !l.dbgMode {
		return
	}
	l.log.Debugln(kv...)
}

func (l *LogrusWrapper) Info(kv ...interface{}) {
	l.log.Infoln(kv...)
}

func (l *LogrusWrapper) Log(kv ...interface{}) {
	l.log.Logln(logrus.InfoLevel, kv...)
}

func (l *LogrusWrapper) Panic(kv ...interface{}) {
	l.log.Panic(kv...)
}

// TODO slog impl
