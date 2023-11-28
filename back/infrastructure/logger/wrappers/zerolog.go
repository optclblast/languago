package wrappers

import (
	"bytes"
	"fmt"

	"github.com/rs/zerolog"
)

type zerologWrapper struct {
	dbgMode bool
	log     zerolog.Logger
}

func (l *zerologWrapper) Write(args ...any) {
	//l.log.Trace().Msgf(format, args...)
}

func (l *zerologWrapper) Warn(args ...any) {
	l.log.Warn().MsgFunc(messageBuilder(args))
}

func (l *zerologWrapper) Debug(args ...any) {
	if !l.dbgMode {
		return
	}
	l.log.Debug().MsgFunc(messageBuilder(args))
}

func (l *zerologWrapper) Error(args ...any) {
	l.log.Error().MsgFunc(messageBuilder(args))
}

func (l *zerologWrapper) Info(args ...any) {
	l.log.Info().MsgFunc(messageBuilder(args))
}

func (l *zerologWrapper) Log(args ...any) {
	//l.log.Log(l.lvl, args) // TODO level mapper
	l.Info(args)
}

func (l *zerologWrapper) Panic(args ...any) {
	l.log.Panic().MsgFunc(messageBuilder(args))
}

func (l *zerologWrapper) Trace(args ...any) {
	l.log.Trace().MsgFunc(messageBuilder(args))
}

// Formatted

func (l *zerologWrapper) Warnf(format string, args ...any) {
	l.log.Warn().Msgf(format, args...)
}

func (l *zerologWrapper) Debugf(format string, args ...any) {
	if !l.dbgMode {
		return
	}
	l.log.Debug().Msgf(format, args...)
}

func (l *zerologWrapper) Errorf(format string, args ...any) {
	l.log.Error().Msgf(format, args...)
}

func (l *zerologWrapper) Infof(format string, args ...any) {
	l.log.Info().Msgf(format, args...)
}

func (l *zerologWrapper) Logf(format string, args ...any) {
	//l.log.Log(l.lvl, args) // TODO level mapper
	l.Infof(format, args...)
}

func (l *zerologWrapper) Panicf(format string, args ...any) {
	l.log.Panic().Msgf(format, args...)
}

func (l *zerologWrapper) Tracef(format string, args ...any) {
	l.log.Trace().Msgf(format, args...)
}

func (l *zerologWrapper) Writef(format string, args ...any) {
	//l.log.Trace().Msgf(format, args...)
}

func messageBuilder(args ...any) func() string {
	return func() string {
		var buf bytes.Buffer

		for _, arg := range args {
			buf.WriteString(fmt.Sprint(arg))
		}

		fmt.Println(buf.String())

		return buf.String()
	}
}
