package logger

import (
	"fmt"
	"time"
)

// Extends std Logger type to fit the Logger interface
func (l *DefaultLogger) Warn(msg string, kv LogFields) {
	l.log.SetPrefix(fmt.Sprintf("[WARN] %v", time.Now()))
	l.log.Println(getFields(msg, kv)...)
}
func (l *DefaultLogger) Err(msg string, kv LogFields) {
	l.log.SetPrefix(fmt.Sprintf("[ERROR] %v", time.Now()))
	l.log.Println(getFields(msg, kv)...)
}
func (l *DefaultLogger) Debug(msg string, kv LogFields) {
	if !l.dbgMode {
		return
	}
	l.log.SetPrefix(fmt.Sprintf("[DEBUG] %v", time.Now()))
	l.log.Println(getFields(msg, kv)...)
}
func (l *DefaultLogger) Info(msg string, kv LogFields) {
	l.log.SetPrefix(fmt.Sprintf("[INFO] %v", time.Now()))
	l.log.Println(getFields(msg, kv)...)
}
func (l *DefaultLogger) Log(msg string, kv LogFields) {
	l.log.SetPrefix(fmt.Sprintf("[LOG] %v", time.Now()))
	l.log.Println(getFields(msg, kv)...)
}
func (l *DefaultLogger) Panic(msg string, kv LogFields) {
	l.log.SetPrefix(fmt.Sprintf("[PANIC] %v", time.Now()))
	l.log.Panicln(getFields(msg, kv)...)
}
