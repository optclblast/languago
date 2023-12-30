package wrappers

import (
	"time"

	"github.com/sirupsen/logrus"
)

type logrusWrapper struct {
	dbgMode bool
	log     *logrus.Logger
}

func (l *logrusWrapper) Write(args ...any) {
	//l.log.Trace(format, args...)
}

func (l *logrusWrapper) Warn(args ...any) {
	l.log.WithTime(time.Now()).Warn(args...)
}

func (l *logrusWrapper) Debug(args ...any) {
	if !l.dbgMode {
		return
	}
	l.log.WithTime(time.Now()).Debug(args...)
}

func (l *logrusWrapper) Error(args ...any) {
	l.log.WithTime(time.Now()).Error(args...)
}

func (l *logrusWrapper) Info(args ...any) {
	l.log.WithTime(time.Now()).Info(args...)
}

func (l *logrusWrapper) Log(args ...any) {
	//l.log.Log(l.lvl, args) // TODO level mapper
	l.Info(args...)
}

func (l *logrusWrapper) Panic(args ...any) {
	l.log.WithTime(time.Now()).Panic(args...)
}

func (l *logrusWrapper) Trace(args ...any) {
	l.log.WithTime(time.Now()).Trace(args...)
}

// Formatted

func (l *logrusWrapper) Warnf(format string, args ...any) {
	l.log.WithTime(time.Now()).Warnf(format, args...)
}

func (l *logrusWrapper) Debugf(format string, args ...any) {
	if !l.dbgMode {
		return
	}
	l.log.WithTime(time.Now()).Debugf(format, args...)
}

func (l *logrusWrapper) Errorf(format string, args ...any) {
	l.log.WithTime(time.Now()).Errorf(format, args...)
}

func (l *logrusWrapper) Infof(format string, args ...any) {
	l.log.WithTime(time.Now()).Infof(format, args...)
}

func (l *logrusWrapper) Logf(format string, args ...any) {
	//l.log.Log(l.lvl, args) // TODO level mapper
	l.Infof(format, args...)
}

func (l *logrusWrapper) Panicf(format string, args ...any) {
	l.log.WithTime(time.Now()).Panicf(format, args...)
}

func (l *logrusWrapper) Tracef(format string, args ...any) {
	l.log.WithTime(time.Now()).Tracef(format, args...)
}

func (l *logrusWrapper) Writef(format string, args ...any) {
	//l.log.Trace(format, args...)
}

func (l *logrusWrapper) do() {
	l.log.WithTime(time.Now()).Info()
}
