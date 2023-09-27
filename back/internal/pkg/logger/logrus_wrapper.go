package logger

import (
	"github.com/sirupsen/logrus"
)

func (l *LogrusWrapper) Warn(msg string, kv LogFields) {
	l.log.Warnln(getFields(msg, kv)...)
}

func (l *LogrusWrapper) Debug(msg string, kv LogFields) {
	if !l.dbgMode {
		return
	}
	l.log.Debugln(getFields(msg, kv)...)
}

func (l *LogrusWrapper) Info(msg string, kv LogFields) {
	l.log.Infoln(getFields(msg, kv)...)
}

func (l *LogrusWrapper) Log(msg string, kv LogFields) {
	l.log.Logln(logrus.InfoLevel, getFields(msg, kv)...)
}

func (l *LogrusWrapper) Panic(msg string, kv LogFields) {
	l.log.Panic(getFields(msg, kv)...)
}

func getFields(msg string, kv LogFields) []interface{} {
	var f []interface{}
	f = append(f, msg)
	for k, v := range kv {
		f = append(f, k)
		f = append(f, v)
	}
	return f
}
