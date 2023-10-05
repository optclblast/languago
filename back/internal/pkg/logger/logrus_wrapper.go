package logger

import (
	"github.com/sirupsen/logrus"
)

func (l *LogrusWrapper) Warn(msg string, kv LogFields) {
	l.log.WithFields(getLogrusFields(kv)).Warnln(msg)
}

func (l *LogrusWrapper) Debug(msg string, kv LogFields) {
	if !l.dbgMode {
		return
	}
	l.log.WithFields(getLogrusFields(kv)).Debugln(msg)
}

func (l *LogrusWrapper) Info(msg string, kv LogFields) {
	l.log.WithFields(getLogrusFields(kv)).Infoln(msg)
}

func (l *LogrusWrapper) Log(msg string, kv LogFields) {
	l.log.WithFields(getLogrusFields(kv)).Log(logrus.DebugLevel, msg)
}

func (l *LogrusWrapper) Panic(msg string, kv LogFields) {
	l.log.WithFields(getLogrusFields(kv)).Panic(msg)
}

func getLogrusFields(kv LogFields) logrus.Fields {
	var fields logrus.Fields = make(logrus.Fields, len(kv))

	for k, v := range kv {
		fields[k] = v
	}
	return fields
}
