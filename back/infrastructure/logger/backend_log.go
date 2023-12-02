package logger

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sync/atomic"

	"languago/pkg/mstime"
)

func (l *logger) Trace(args ...interface{}) {
	l.Write(LevelTrace, args...)
}

func (l *logger) Tracef(format string, args ...interface{}) {
	l.Writef(LevelTrace, format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.Write(LevelDebug, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.Writef(LevelDebug, format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.Write(LevelInfo, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Writef(LevelInfo, format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.Write(LevelWarn, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Writef(LevelWarn, format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.Write(LevelError, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Writef(LevelError, format, args...)
}

func (l *logger) Critical(args ...interface{}) {
	l.Write(LevelCritical, args...)
}

func (l *logger) Criticalf(format string, args ...interface{}) {
	l.Writef(LevelCritical, format, args...)
}

func (l *logger) Write(logLevel Level, args ...interface{}) {
	lvl := l.Level()
	if lvl <= logLevel {
		l.print(logLevel, l.tag, args...)
	}
}

func (l *logger) Writef(logLevel Level, format string, args ...interface{}) {
	lvl := l.Level()
	if lvl <= logLevel {
		l.printf(logLevel, l.tag, format, args...)
	}
}

func (l *logger) Level() Level {
	return Level(atomic.LoadUint32((*uint32)(&l.lvl)))
}

func (l *logger) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&l.lvl), uint32(level))
}

func (l *logger) Backend() *Backend {
	return l.b
}

func (l *logger) printf(lvl Level, tag string, format string, args ...interface{}) {
	t := mstime.Now() // get as early as possible

	var file string
	var line int
	if l.b.flag&(LogFlagShortFile|LogFlagLongFile) != 0 {
		file, line = callsite(l.b.flag)
	}

	buf := make([]byte, 0, normalLogSize)

	formatHeader(&buf, t, lvl.String(), tag, file, line)
	bytesBuf := bytes.NewBuffer(buf)
	_, _ = fmt.Fprintf(bytesBuf, format, args...)
	bytesBuf.WriteByte('\n')

	if !l.b.IsRunning() {
		_, _ = fmt.Fprintf(os.Stderr, bytesBuf.String())
		panic("Writing to the logger when it's not running")
	}
	l.writeChan <- logEntry{bytesBuf.Bytes(), lvl}
}

func (l *logger) print(lvl Level, tag string, args ...interface{}) {
	if atomic.LoadUint32(&l.b.isRunning) == 0 {
		panic("printing log without initializing")
	}
	t := mstime.Now() // get as early as possible

	var file string
	var line int
	if l.b.flag&(LogFlagShortFile|LogFlagLongFile) != 0 {
		file, line = callsite(l.b.flag)
	}

	buf := make([]byte, 0, normalLogSize)
	formatHeader(&buf, t, lvl.String(), tag, file, line)
	bytesBuf := bytes.NewBuffer(buf)
	_, _ = fmt.Fprintln(bytesBuf, args...)

	if !l.b.IsRunning() {
		panic("Writing to the logger when it's not running")
	}
	l.writeChan <- logEntry{bytesBuf.Bytes(), lvl}
}

// From stdlib log package.
// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func formatHeader(buf *[]byte, t mstime.Time, lvl, tag string, file string, line int) {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	ms := t.Millisecond()

	itoa(buf, year, 4)
	*buf = append(*buf, '-')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '-')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)
	*buf = append(*buf, '.')
	itoa(buf, ms, 3)
	*buf = append(*buf, " ["...)
	*buf = append(*buf, lvl...)
	*buf = append(*buf, "] "...)
	*buf = append(*buf, tag...)
	if file != "" {
		*buf = append(*buf, ' ')
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
	}
	*buf = append(*buf, ": "...)
}

const calldepth = 4

func callsite(flag uint32) (string, int) {
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		return "???", 0
	}
	if flag&LogFlagShortFile != 0 {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if os.IsPathSeparator(file[i]) {
				short = file[i+1:]
				break
			}
		}
		file = short
	}
	return file, line
}
