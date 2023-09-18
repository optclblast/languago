package config

import (
	"log"

	"log/slog"

	"github.com/sirupsen/logrus"
)

type (
	Logger interface {
		Warn(kv ...interface{})
		Debug(kv ...interface{})
		Info(kv ...interface{})
		Log(kv ...interface{})
	}

	DefaultLogger struct {
		dbgMode bool
		log     *log.Logger
	}

	SlogLogger struct {
		dbgMode bool
		log     *slog.Logger
	}

	LogrusWrapper struct {
		dbgMode bool
		log     *logrus.Entry
	}
)

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

func (l *DefaultLogger) Warn(kv ...interface{}) {
	l.log.Warn(kv...)
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
