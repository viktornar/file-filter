package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	ErrorLevel
	DisabledLevel
)

type ExternalLogger interface {
	Log(Level, string)
	Flush()
}

var (
	Debug = &logger{DebugLevel}
	Info  = &logger{InfoLevel}
	Error = &logger{ErrorLevel}
)

type globalState struct {
	currentLevel  Level
	defaultLogger Logger
	external      ExternalLogger
}

var (
	mu    sync.RWMutex
	state = globalState{
		currentLevel:  InfoLevel,
		defaultLogger: newDefaultLogger(os.Stderr),
	}
)

func globals() globalState {
	mu.RLock()
	defer mu.RUnlock()
	return state
}

func newDefaultLogger(w io.Writer) Logger {
	return log.New(w, "", log.Ldate|log.Ltime|log.LUTC|log.Lmicroseconds)
}

type logBridge struct {
	Logger
}

func (lb logBridge) Write(b []byte) (n int, err error) {
	var message string
	// Split "f.go:42: message" into "f.go", "42", and "message".
	parts := bytes.SplitN(b, []byte{':'}, 3)
	if len(parts) != 3 || len(parts[0]) < 1 || len(parts[2]) < 1 {
		message = fmt.Sprintf("bad log format: %s", b)
	} else {
		message = string(parts[2][1:]) // Skip leading space.
	}
	lb.Print(message)
	return len(b), nil
}

// NewStdLogger creates a *log.Logger ("log" is from the Go standard library)
// that forwards messages to the provided upspin logger using a logBridge. The
// standard logger is configured with log.Lshortfile, this log line
// format which is parsed to extract the log message (skipping the filename,
// line number) to forward it to the provided upspin logger.
func NewStdLogger(l Logger) *log.Logger {
	lb := logBridge{l}
	return log.New(lb, "", log.Lshortfile)
}

// Register connects an ExternalLogger to the default logger. This may only be
// called once.
func Register(e ExternalLogger) {
	mu.Lock()
	defer mu.Unlock()

	if state.external != nil {
		panic("cannot register second external logger")
	}
	state.external = e
}

// SetOutput sets the default loggers to write to w.
// If w is nil, the default loggers are disabled.
func SetOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()

	if w == nil {
		state.defaultLogger = nil
	} else {
		state.defaultLogger = newDefaultLogger(w)
	}
}

type logger struct {
	level Level
}

var _ Logger = (*logger)(nil)

func (l *logger) Printf(format string, v ...interface{}) {
	g := globals()

	if l.level < g.currentLevel {
		return
	}

	if g.external != nil {
		g.external.Log(l.level, fmt.Sprintf(format, v...))
	}

	if g.defaultLogger != nil {
		g.defaultLogger.Printf(format, v...)
	}
}

// Print writes a message to the log.
func (l *logger) Print(v ...interface{}) {
	g := globals()

	if l.level < g.currentLevel {
		return // Don't log at lower levels.
	}

	if g.external != nil {
		g.external.Log(l.level, fmt.Sprint(v...))
	}

	if g.defaultLogger != nil {
		g.defaultLogger.Print(v...)
	}
}

func (l *logger) Println(v ...interface{}) {
	g := globals()

	if l.level < g.currentLevel {
		return // Don't log at lower levels.
	}
	if g.external != nil {
		g.external.Log(l.level, fmt.Sprintln(v...))
	}
	if g.defaultLogger != nil {
		g.defaultLogger.Println(v...)
	}
}

func (l *logger) Fatal(v ...interface{}) {
	g := globals()

	if g.external != nil {
		g.external.Log(l.level, fmt.Sprint(v...))
		// Make sure we get the Fatal recorded.
		g.external.Flush()
		// Fall through to ensure we record it locally too.
	}
	if g.defaultLogger != nil {
		g.defaultLogger.Fatal(v...)
	} else {
		log.Fatal(v...)
	}
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	g := globals()

	if g.external != nil {
		g.external.Log(l.level, fmt.Sprintf(format, v...))
		// Make sure we get the Fatal recorded.
		g.external.Flush()
		// Fall through to ensure we record it locally too.
	}
	if g.defaultLogger != nil {
		g.defaultLogger.Fatalf(format, v...)
	} else {
		log.Fatalf(format, v...)
	}
}

func (l *logger) Flush() {
	Flush()
}

func (l *logger) String() string {
	return toString(l.level)
}

func toString(level Level) string {
	switch level {
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	case ErrorLevel:
		return "error"
	case DisabledLevel:
		return "disabled"
	}
	return "unknown"
}

func toLevel(level string) (Level, error) {
	switch level {
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	case "error":
		return ErrorLevel, nil
	case "disabled":
		return DisabledLevel, nil
	}
	return DisabledLevel, fmt.Errorf("invalid log level %q", level)
}

func GetLevel() string {
	g := globals()

	return toString(g.currentLevel)
}

func SetLevel(level string) error {
	l, err := toLevel(level)
	if err != nil {
		return err
	}
	mu.Lock()
	state.currentLevel = l
	mu.Unlock()
	return nil
}

func At(level string) bool {
	g := globals()

	l, err := toLevel(level)
	if err != nil {
		return false
	}
	return g.currentLevel <= l
}

func Printf(format string, v ...interface{}) {
	Info.Printf(format, v...)
}

func Print(v ...interface{}) {
	Info.Print(v...)
}

func Println(v ...interface{}) {
	Info.Println(v...)
}

func Fatal(v ...interface{}) {
	Info.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	Info.Fatalf(format, v...)
}

func Flush() {
	g := globals()

	if g.external != nil {
		g.external.Flush()
	}
}
