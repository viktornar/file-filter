package logger

import (
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

var (
	Debug = &logger{DebugLevel}
	Info  = &logger{InfoLevel}
	Error = &logger{ErrorLevel}
)

type globalState struct {
	currentLevel  Level
	defaultLogger Logger
}

var (
	mu    sync.RWMutex
	state = globalState{
		currentLevel:  InfoLevel,
		defaultLogger: newDefaultLogger(os.Stdout),
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

// SetOutput sets the logger to write to w.
// If w is nil, the logger is just disabled.
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

// TODO: Add log level at start of the line
func (l *logger) Printf(format string, v ...any) {
	g := globals()

	if l.level < g.currentLevel {
		return
	}

	if g.defaultLogger != nil {
		g.defaultLogger.Printf(format, v...)
	}
}

func (l *logger) Print(v ...any) {
	g := globals()

	if l.level < g.currentLevel {
		return // Don't log at lower levels.
	}

	if g.defaultLogger != nil {
		g.defaultLogger.Print(v...)
	}
}

func (l *logger) Println(v ...any) {
	g := globals()

	if l.level < g.currentLevel {
		return // Don't log at lower levels.
	}

	if g.defaultLogger != nil {
		g.defaultLogger.Println(v...)
	}
}

func (l *logger) Fatal(v ...any) {
	g := globals()

	if g.defaultLogger != nil {
		g.defaultLogger.Fatal(v...)
	} else {
		log.Fatal(v...)
	}
}

func (l *logger) Fatalf(format string, v ...any) {
	g := globals()

	if g.defaultLogger != nil {
		g.defaultLogger.Fatalf(format, v...)
	} else {
		log.Fatalf(format, v...)
	}
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

func Printf(format string, v ...any) {
	Info.Printf(format, v...)
}

func Print(v ...any) {
	Info.Print(v...)
}

func Println(v ...any) {
	Info.Println(v...)
}

func Fatal(v ...any) {
	Info.Fatal(v...)
}

func Fatalf(format string, v ...any) {
	Info.Fatalf(format, v...)
}
