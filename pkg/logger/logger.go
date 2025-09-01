package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// represents the severity of the log
type LogLevel int

const (
	INFO LogLevel = iota
	SUCCESS
	WARN
	ERROR
)

// prints colored and timestamped messages
type Logger struct{}

// creates a new Logger instance
func NewLogger() *Logger {
	return &Logger{}
}

// Logger instance
var Log = NewLogger()

// Log prints a message with the given level
func (l *Logger) Log(level LogLevel, format string, args ...any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)

	color.New(color.FgHiBlack).Printf("[%s] ", timestamp)

	switch level {
	case INFO:
		color.New(color.FgCyan).Printf("[INFO] ")
	case SUCCESS:
		color.New(color.FgGreen).Printf("[SUCCESS] ")
	case WARN:
		color.New(color.FgYellow).Printf("[WARN] ")
	case ERROR:
		color.New(color.FgRed).Printf("[ERROR] ")
	default:
		fmt.Printf("[UNKNOWN] ")
	}

	color.New(color.FgWhite).Printf("%s\n", msg)
}

// prints an info message
func (l *Logger) Info(format string, args ...any) {
	l.Log(INFO, format, args...)
}

// prints a success message
func (l *Logger) Success(format string, args ...any) {
	l.Log(SUCCESS, format, args...)
}

// prints a warning message
func (l *Logger) Warn(format string, args ...any) {
	l.Log(WARN, format, args...)
}

// prints an error message
func (l *Logger) Error(format string, args ...any) {
	l.Log(ERROR, format, args...)
}
