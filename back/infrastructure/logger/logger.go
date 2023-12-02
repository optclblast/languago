package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type (
	// todo just use logrus or zerolog

	abstractLoggerConfig interface {
		GetLevel() Level
		GetEnv() EnvParam

		//deprecated
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

type EnvParam string

const (
	EnvParam_LOCAL       EnvParam = "local"
	EnvParam_DEVELOPMENT EnvParam = "development"
	EnvParam_PRODUCTION  EnvParam = "production"
)

func ProvideLogger(cfg abstractLoggerConfig) zerolog.Logger {
	switch cfg.GetEnv() {
	case EnvParam_LOCAL:
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
		return zerolog.New(consoleWriter).With().Timestamp().Logger().Level(zerolog.Level(cfg.GetLevel()))
	case EnvParam_DEVELOPMENT:
		// todo
	case EnvParam_PRODUCTION:
		// todo
	default:
		fmt.Fprintln(os.Stdout, "fatal error: invalid env parameter")
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	return zerolog.New(consoleWriter).With().Timestamp().Logger().Level(zerolog.Level(cfg.GetLevel()))
}
