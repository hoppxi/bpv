package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

// represents the severity of the log
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	SUCCESS
	WARN
	ERROR
	FATAL
)

// Logger configuration
type Logger struct {
	verbose bool
}

// creates a new Logger instance
func NewLogger(verbose bool) *Logger {
	return &Logger{verbose: verbose}
}

// Global logger instance
var Log *Logger

// Initialize the global logger
func Init(verbose bool) {
	Log = NewLogger(verbose)
}

// Log prints a message with the given level
func (l *Logger) Log(level LogLevel, format string, args ...any) {
	if (level == DEBUG || level == INFO) && (!l.verbose) {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)

	color.New(color.FgHiBlack).Printf("[%s] ", timestamp)

	switch level {
	case DEBUG:
		color.New(color.FgMagenta).Printf("[DEBUG] ")
	case INFO:
		color.New(color.FgCyan).Printf("[INFO] ")
	case SUCCESS:
		color.New(color.FgGreen).Printf("[SUCCESS] ")
	case WARN:
		color.New(color.FgYellow).Printf("[WARN] ")
	case ERROR:
		color.New(color.FgRed).Printf("[ERROR] ")
	case FATAL:
		color.New(color.FgRed, color.Bold).Printf("[FATAL] ")
	default:
		fmt.Printf("[UNKNOWN] ")
	}

	color.New(color.FgWhite).Printf("%s\n", msg)
}

// Debug prints a debug message (only in verbose mode)
func (l *Logger) Debug(format string, args ...any) {
	l.Log(DEBUG, format, args...)
}

// Info prints an info message
func (l *Logger) Info(format string, args ...any) {
	l.Log(INFO, format, args...)
}

// Success prints a success message
func (l *Logger) Success(format string, args ...any) {
	l.Log(SUCCESS, format, args...)
}

// Warn prints a warning message
func (l *Logger) Warn(format string, args ...any) {
	l.Log(WARN, format, args...)
}

// Error prints an error message and returns formatted error
func (l *Logger) Error(format string, args ...any) error {
	errMsg := fmt.Sprintf(format, args...)
	l.Log(ERROR, errMsg)
	return fmt.Errorf(errMsg)
}

// ErrorP prints an error message with prefix and returns formatted error
func (l *Logger) ErrorP(prefix, format string, args ...any) error {
	errMsg := fmt.Sprintf("%s: %s", prefix, fmt.Sprintf(format, args...))
	l.Log(ERROR, errMsg)
	return fmt.Errorf(errMsg)
}

// Fatal prints a fatal error message and exits the program
func (l *Logger) Fatal(format string, args ...any) {
	l.Log(FATAL, format, args...)
	os.Exit(1)
}

// FatalP prints a fatal error message with prefix and exits the program
func (l *Logger) FatalP(prefix, format string, args ...any) {
	errMsg := fmt.Sprintf("%s: %s", prefix, fmt.Sprintf(format, args...))
	l.Log(FATAL, errMsg)
	os.Exit(1)
}

// FatalErr prints a fatal error from an existing error and exits
func (l *Logger) FatalErr(err error, context ...string) {
	if len(context) > 0 {
		l.Log(FATAL, "%s: %v", context[0], err)
	} else {
		l.Log(FATAL, "%v", err)
	}
	os.Exit(1)
}

// IsVerbose returns whether verbose logging is enabled
func (l *Logger) IsVerbose() bool {
	return l.verbose
}

// LogDuration logs the duration of an operation
func (l *Logger) LogDuration(start time.Time, operation string) {
	duration := time.Since(start)
	l.Debug("%s completed in %v", operation, duration)
}

// LogIfError logs an error if it exists and returns it
func (l *Logger) LogIfError(err error, context string) error {
	if err != nil {
		l.ErrorP(context, "%v", err)
	}
	return err
}

// LogIfErrorP logs an error with custom message if it exists and returns it
func (l *Logger) LogIfErrorP(err error, format string, args ...any) error {
	if err != nil {
		l.ErrorP(fmt.Sprintf(format, args...), "%v", err)
	}
	return err
}