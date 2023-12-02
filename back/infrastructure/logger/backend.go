package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"

	"errors"

	"github.com/jrick/logrotate/rotator"
)

const logsBuffer = 0

const normalLogSize = 512

var defaultFlags = getDefaultFlags()

func getDefaultFlags() (flags uint32) {
	for _, f := range strings.Split(os.Getenv("LOGFLAGS"), ",") {
		switch f {
		case "longfile":
			flags |= LogFlagLongFile
		case "shortfile":
			flags |= LogFlagShortFile
		}
	}
	return
}

// Flags to modify Backend's behavior.
const (
	LogFlagLongFile uint32 = 1 << iota
	LogFlagShortFile
)

type Backend struct {
	flag      uint32
	isRunning uint32
	writers   []logWriter
	writeChan chan logEntry
	syncClose sync.Mutex // used to sync that the logger finished writing everything
}

func NewBackendWithFlags(flags uint32) *Backend {
	return &Backend{flag: flags, writeChan: make(chan logEntry, logsBuffer)}
}

func NewBackend() *Backend {
	return NewBackendWithFlags(defaultFlags)
}

const (
	defaultThresholdKB = 100 * 1000 // 100 MB logs by default.
	defaultMaxRolls    = 8          // keep 8 last logs by default.
)

type logWriter interface {
	io.WriteCloser
	LogLevel() Level
}

type logWriterWrap struct {
	io.WriteCloser
	logLevel Level
}

func (lw logWriterWrap) LogLevel() Level {
	return lw.logLevel
}

func (b *Backend) AddLogFile(logFile string, logLevel Level) error {
	return b.AddLogFileWithCustomRotator(logFile, logLevel, defaultThresholdKB, defaultMaxRolls)
}

func (b *Backend) AddLogWriter(logWriter io.WriteCloser, logLevel Level) error {
	if b.IsRunning() {
		return errors.New("The logger is already running")
	}
	b.writers = append(b.writers, logWriterWrap{
		WriteCloser: logWriter,
		logLevel:    logLevel,
	})
	return nil
}

func (b *Backend) AddLogFileWithCustomRotator(logFile string, logLevel Level, thresholdKB int64, maxRolls int) error {
	if b.IsRunning() {
		return errors.New("The logger is already running")
	}
	logDir, _ := filepath.Split(logFile)
	if logDir != "" {
		err := os.MkdirAll(logDir, 0700)
		if err != nil {
			return fmt.Errorf("failed to create log directory: %+v", err)
		}
	}
	r, err := rotator.New(logFile, thresholdKB, false, maxRolls)
	if err != nil {
		return fmt.Errorf("failed to create file rotator: %s", err)
	}
	b.writers = append(b.writers, logWriterWrap{
		WriteCloser: r,
		logLevel:    logLevel,
	})
	return nil
}

func (b *Backend) Run() error {
	if !atomic.CompareAndSwapUint32(&b.isRunning, 0, 1) {
		return errors.New("The logger is already running")
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Fatal error in logger.Backend goroutine: %+v\n", err)
				_, _ = fmt.Fprintf(os.Stderr, "Goroutine stacktrace: %s\n", debug.Stack())
			}
		}()
		b.runBlocking()
	}()
	return nil
}

func (b *Backend) runBlocking() {
	defer atomic.StoreUint32(&b.isRunning, 0)
	b.syncClose.Lock()
	defer b.syncClose.Unlock()

	for log := range b.writeChan {
		for _, writer := range b.writers {
			if log.level >= writer.LogLevel() {
				_, _ = writer.Write(log.log)
			}
		}
	}
}

func (b *Backend) IsRunning() bool {
	return atomic.LoadUint32(&b.isRunning) != 0
}

func (b *Backend) Close() {
	close(b.writeChan)
	b.syncClose.Lock()
	defer b.syncClose.Unlock()
	for _, writer := range b.writers {
		_ = writer.Close()
	}
}

// func (b *Backend) Logger(subsystemTag string) *Logger {
// 	return &Logger{LevelOff, subsystemTag, b, b.writeChan}
// }
