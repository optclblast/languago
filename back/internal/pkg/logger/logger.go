package logger

import (
	"log"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
		log     *logrus.Logger
	}

	Logger interface {
		Warn(msg string, kv LogFields)
		Error(msg string, kv LogFields)
		Debug(msg string, kv LogFields)
		Info(msg string, kv LogFields)
		Log(msg string, kv LogFields)
		Panic(msg string, kv LogFields)
	}

	LogFields map[string]interface{}

	EnvParam string
)

const (
	MessageField string = "message"
	ErrorField   string = "error"
	ContentField string = "content"
	////////////////////////////////
	EnvParam_LOCAL       EnvParam = "local"
	EnvParam_DEVELOPMENT EnvParam = "development"
	EnvParam_PRODUCTION  EnvParam = "production"
)

func MustToEnvParam(raw string) EnvParam {
	switch {
	case raw == "local":
		return EnvParam_LOCAL
	case raw == "development":
		return EnvParam_DEVELOPMENT
	case raw == "production":
		return EnvParam_PRODUCTION
	default:
		log.Fatalln("fatal error: invalid env parameter")
	}
	return EnvParam_LOCAL
}

func NewLogrusWrapper(dbg bool, env EnvParam) *LogrusWrapper {
	var llog *logrus.Logger

	// TODO graylog output
	// llog.SetOutput()
	switch env {
	case EnvParam_LOCAL:
		llog = logrus.StandardLogger()
		llog.Formatter = &logrus.TextFormatter{
			ForceColors:      true,
			QuoteEmptyFields: true,
		}
		llog.SetLevel(logrus.DebugLevel)
	case EnvParam_DEVELOPMENT:
		llog = logrus.New()
		llog.Formatter = &logrus.JSONFormatter{
			DisableHTMLEscape: true,
			PrettyPrint:       true,
		}
		llog.SetLevel(logrus.DebugLevel)
	case EnvParam_PRODUCTION:
		llog = logrus.New()
		llog.Formatter = &logrus.JSONFormatter{
			DisableHTMLEscape: true,
		}
		llog.SetLevel(logrus.InfoLevel)
	default:
		log.Fatalln("fatal error: invalid env parameter")
	}

	lw := &LogrusWrapper{
		dbgMode: dbg,
		log:     llog,
	}
	return lw
}

func NewZapWrapper(dbg bool, env EnvParam) *ZapWrapper {
	var (
		opts []zap.Option = make([]zap.Option, 0)
		zlog *zap.Logger
	)

	switch env {
	case EnvParam_LOCAL:
		opts = append(opts,
			zap.Development(),
			zap.WithCaller(true),
			zap.WithClock(zapcore.DefaultClock),
		)
		zlog = zap.Must(zap.NewDevelopment(opts...))
	case EnvParam_DEVELOPMENT:
		opts = append(opts,
			zap.Development(),
			zap.WithCaller(true),
			zap.WithClock(zapcore.DefaultClock),
		)
		zlog = zap.Must(zap.NewProduction(opts...))
	case EnvParam_PRODUCTION:
		opts = append(opts,
			zap.WithClock(zapcore.DefaultClock),
			zap.IncreaseLevel(zapcore.InfoLevel),
		)
		zlog = zap.Must(zap.NewProduction(opts...))
	default:
		log.Fatalln("fatal error: invalid env parameter")
	}

	return &ZapWrapper{
		dbgMode: dbg,
		log:     zlog,
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

func provideDefaultLogger(dbg bool) *DefaultLogger {
	logger := DefaultLogger{
		dbgMode: dbg,
		log:     log.Default(),
	}
	return &logger
}

func LogFieldPair(key string, val any) LogFields {
	lf := make(LogFields, 0)
	lf[key] = val
	return lf
}
