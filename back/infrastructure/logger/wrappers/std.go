package wrappers

import "log"

// TODO stdlog --> slog

type DefaultLogger struct {
	dbgMode bool
	log     *log.Logger
}

// Extends std Logger type to fit the Logger interface
// func (l *DefaultLogger) Warn(args ...any) {
// 	l.log.SetPrefix(fmt.Sprintf("[WARN] %v", time.Now()))
// 	l.log.Println(getStdFields(msg, kv)...)
// }
// func (l *DefaultLogger) Error(args ...any) {
// 	l.log.SetPrefix(fmt.Sprintf("[ERROR] %v", time.Now()))
// 	l.log.Println(getStdFields(msg, kv)...)
// }
// func (l *DefaultLogger) Debug(args ...any) {
// 	if !l.dbgMode {
// 		return
// 	}
// 	l.log.SetPrefix(fmt.Sprintf("[DEBUG] %v", time.Now()))
// 	l.log.Println(getStdFields(msg, kv)...)
// }
// func (l *DefaultLogger) Info(args ...any) {
// 	l.log.SetPrefix(fmt.Sprintf("[INFO] %v", time.Now()))
// 	l.log.Println(getStdFields(msg, kv)...)
// }
// func (l *DefaultLogger) Log(args ...any) {
// 	l.log.SetPrefix(fmt.Sprintf("[LOG] %v", time.Now()))
// 	l.log.Println(getStdFields(msg, kv)...)
// }
// func (l *DefaultLogger) Panic(args ...any) {
// 	l.log.SetPrefix(fmt.Sprintf("[PANIC] %v", time.Now()))
// 	l.log.Panicln(getStdFields(msg, kv)...)
// }
