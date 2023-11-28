package wrappers

import (
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
)

type EnvParam string

const (
	MessageField string = "message"
	ErrorField   string = "error"
	ContentField string = "content"

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

func NewLogrusWrapper(dbg bool, env EnvParam) *logrusWrapper {
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

	lw := &logrusWrapper{
		dbgMode: dbg,
		log:     llog,
	}
	return lw
}

func NewZerologWrapper(dbg bool, env EnvParam) *zerologWrapper {
	var (
		logger zerolog.Logger
	)

	switch env {
	case EnvParam_LOCAL:
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
		//multi := zerolog.MultiLevelWriter(consoleWriter, os.Stdout)
		logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
	case EnvParam_DEVELOPMENT:
		// logger = zerolog.New()
	case EnvParam_PRODUCTION:
		//
	default:
		log.Fatalln("fatal error: invalid env parameter")
	}

	return &zerologWrapper{
		dbgMode: dbg,
		log:     logger,
	}
}
