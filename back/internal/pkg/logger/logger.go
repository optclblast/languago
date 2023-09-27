package logger

import (
	"log"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type (
	abstractLoggerConfig interface {
		GetLogger() Logger
	}

	DefaultLogger struct {
		dbgMode bool
		log     *log.Logger
	}

	ZapWrapper struct {
		dbgMode bool
		log     *zap.Logger
	}

	LogrusWrapper struct {
		dbgMode bool
		log     *logrus.Entry
	}

	Logger interface {
		Warn(msg string, kv LogFields)
		//Err(msg string, kv LogFields)
		Debug(msg string, kv LogFields)
		Info(msg string, kv LogFields)
		Log(msg string, kv LogFields)
	}

	LogFields map[string]interface{}
)

const (
	MessageField string = "message"
	ErrorField   string = "error"
	ContentField string = "content"
)

func NewLogrusWrapper(dbg bool) *LogrusWrapper {
	return &LogrusWrapper{
		dbgMode: dbg,
		log:     logrus.NewEntry(logrus.New()),
	}
}

func NewZapWrapper(dbg bool) *ZapWrapper {
	return &ZapWrapper{
		dbgMode: dbg,
		log:     zap.Must(zap.NewProduction()),
	}
}

func NewDefaultLogger(dbg bool) *DefaultLogger {
	return provideDefaultLogger(dbg)
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

func LogFieldPair(key string, val any) LogFields {
	lf := make(LogFields, 0)
	lf[key] = val
	return lf
}
