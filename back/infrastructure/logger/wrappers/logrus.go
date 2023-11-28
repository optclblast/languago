package wrappers

import "github.com/sirupsen/logrus"

type logrusWrapper struct {
	dbgMode bool
	log     *logrus.Logger
}

func (l *logrusWrapper) Write(args ...any) {
	//l.log.Trace().Msgf(format, args...)
}

func (l *logrusWrapper) Warn(args ...any) {
	l.log.Warn(args...)
}

func (l *logrusWrapper) Debug(args ...any) {
	if !l.dbgMode {
		return
	}
	l.log.Debug(args...)
}

func (l *logrusWrapper) Error(args ...any) {
	l.log.Error(args...)
}

func (l *logrusWrapper) Info(args ...any) {
	l.log.Info(args...)
}

func (l *logrusWrapper) Log(args ...any) {
	//l.log.Log(l.lvl, args) // TODO level mapper
	l.Info(args...)
}

func (l *logrusWrapper) Panic(args ...any) {
	l.log.Panic(args...)
}

func (l *logrusWrapper) Trace(args ...any) {
	l.log.Trace(args...)
}

// Formatted

func (l *logrusWrapper) Warnf(format string, args ...any) {
	l.log.Warnf(format, args...)
}

func (l *logrusWrapper) Debugf(format string, args ...any) {
	if !l.dbgMode {
		return
	}
	l.log.Debugf(format, args...)
}

func (l *logrusWrapper) Errorf(format string, args ...any) {
	l.log.Errorf(format, args...)
}

func (l *logrusWrapper) Infof(format string, args ...any) {
	l.log.Infof(format, args...)
}

func (l *logrusWrapper) Logf(format string, args ...any) {
	//l.log.Log(l.lvl, args) // TODO level mapper
	l.Infof(format, args...)
}

func (l *logrusWrapper) Panicf(format string, args ...any) {
	l.log.Panicf(format, args...)
}

func (l *logrusWrapper) Tracef(format string, args ...any) {
	l.log.Tracef(format, args...)
}

func (l *logrusWrapper) Writef(format string, args ...any) {
	//l.log.Trace().Msgf(format, args...)
}
