package logger

import (
	"go.uber.org/zap"
)

func (l *ZapWrapper) Warn(msg string, kv LogFields) {
	l.log.Warn(msg, getZapFields(kv)...)
}

func (l *ZapWrapper) Debug(msg string, kv LogFields) {
	if !l.dbgMode {
		return
	}
	l.log.Debug(msg, getZapFields(kv)...)
}

func (l *ZapWrapper) Info(msg string, kv LogFields) {
	l.log.Info(msg, getZapFields(kv)...)
}

func (l *ZapWrapper) Log(msg string, kv LogFields) {
	l.log.Log(zap.InfoLevel, msg, getZapFields(kv)...)
}

func (l *ZapWrapper) Panic(msg string, kv LogFields) {
	l.log.Panic(msg, getZapFields(kv)...)
}

func getZapFields(kv LogFields) []zap.Field {
	var f []zap.Field
	for k, v := range kv {
		switch v.(type) {
		case string:
			f = append(f, zap.String(k, v.(string)))
		//int
		case int:
			f = append(f, zap.Int(k, v.(int)))
		case int8:
			f = append(f, zap.Int8(k, v.(int8)))
		case int16:
			f = append(f, zap.Int16(k, v.(int16)))
		case int32:
			f = append(f, zap.Int32(k, v.(int32)))
		case int64:
			f = append(f, zap.Int64(k, v.(int64)))
		//uint
		case uint:
			f = append(f, zap.Uint(k, v.(uint)))
		case uint8:
			f = append(f, zap.Uint8(k, v.(uint8)))
		case uint16:
			f = append(f, zap.Uint16(k, v.(uint16)))
		case uint32:
			f = append(f, zap.Uint32(k, v.(uint32)))
		case uint64:
			f = append(f, zap.Uint64(k, v.(uint64)))
		//float
		case float32:
			f = append(f, zap.Float32(k, v.(float32)))
		case float64:
			f = append(f, zap.Float64(k, v.(float64)))
		case bool:
			f = append(f, zap.Bool(k, v.(bool)))
		}
	}
	return f
}
