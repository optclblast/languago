package logger

import "strings"

type Level uint32

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
	LevelOff
)

var levelStrs = [...]string{"TRC", "DBG", "INF", "WRN", "ERR", "CRT", "OFF"}

func LevelFromString(s string) (l Level, ok bool) {
	switch strings.ToLower(s) {
	case "trace", "trc":
		return LevelTrace, true
	case "debug", "dbg":
		return LevelDebug, true
	case "info", "inf":
		return LevelInfo, true
	case "warn", "wrn":
		return LevelWarn, true
	case "error", "err":
		return LevelError, true
	case "critical", "crt":
		return LevelCritical, true
	case "off":
		return LevelOff, true
	default:
		return LevelInfo, false
	}
}

func (l Level) String() string {
	if l >= LevelOff {
		return "OFF"
	}
	return levelStrs[l]
}
