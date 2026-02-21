package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	SUCCESS
	WARN
	ERROR
	FATAL
)

type Logger struct {
	verbose      bool
	colorEnabled bool
	output       io.Writer
}

func NewLogger(verbose, colorEnabled bool) *Logger {
	return &Logger{
		verbose:      verbose,
		colorEnabled: colorEnabled,
		output:       os.Stderr,
	}
}

var Log *Logger

func Init(verbose, colorEnabled bool) {
	Log = NewLogger(verbose, colorEnabled)
}

func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

func (l *Logger) log(level LogLevel, format string, args ...any) {
	if level == DEBUG && !l.verbose {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)

	if l.colorEnabled {
		color.New(color.FgHiBlack).Fprintf(l.output, "[%s] ", timestamp)

		switch level {
		case DEBUG:
			color.New(color.FgMagenta).Fprintf(l.output, "[DEBUG] ")
		case INFO:
			color.New(color.FgCyan).Fprintf(l.output, "[INFO] ")
		case SUCCESS:
			color.New(color.FgGreen).Fprintf(l.output, "[SUCCESS] ")
		case WARN:
			color.New(color.FgYellow).Fprintf(l.output, "[WARN] ")
		case ERROR:
			color.New(color.FgRed).Fprintf(l.output, "[ERROR] ")
		case FATAL:
			color.New(color.FgRed, color.Bold).Fprintf(l.output, "[FATAL] ")
		default:
			fmt.Fprintf(l.output, "[UNKNOWN] ")
		}

		color.New(color.FgWhite).Fprintf(l.output, "%s\n", msg)
	} else {
		levelStr := "UNKNOWN"
		switch level {
		case DEBUG:
			levelStr = "DEBUG"
		case INFO:
			levelStr = "INFO"
		case SUCCESS:
			levelStr = "SUCCESS"
		case WARN:
			levelStr = "WARN"
		case ERROR:
			levelStr = "ERROR"
		case FATAL:
			levelStr = "FATAL"
		}
		fmt.Fprintf(l.output, "[%s] [%s] %s\n", timestamp, levelStr, msg)
	}
}

func (l *Logger) Debug(format string, args ...any) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...any) {
	l.log(INFO, format, args...)
}

func (l *Logger) Success(format string, args ...any) {
	l.log(SUCCESS, format, args...)
}

func (l *Logger) Warn(format string, args ...any) {
	l.log(WARN, format, args...)
}

func (l *Logger) Error(format string, args ...any) error {
	l.log(ERROR, format, args...)
	return fmt.Errorf(format, args...)
}

func (l *Logger) ErrorP(prefix, format string, args ...any) error {
	l.log(ERROR, "%s: %s", prefix, fmt.Sprintf(format, args...))
	return fmt.Errorf("%s: %s", prefix, fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(format string, args ...any) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

func (l *Logger) FatalP(prefix, format string, args ...any) {
	l.log(FATAL, "%s: %s", prefix, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *Logger) FatalErr(err error, context ...string) {
	if len(context) > 0 {
		l.log(FATAL, "%s: %v", context[0], err)
	} else {
		l.log(FATAL, "%v", err)
	}
	os.Exit(1)
}

func (l *Logger) IsVerbose() bool {
	return l.verbose
}

func (l *Logger) LogDuration(start time.Time, operation string) {
	duration := time.Since(start)
	l.Debug("%s completed in %v", operation, duration)
}

func (l *Logger) LogIfError(err error, context string) error {
	if err != nil {
		l.ErrorP(context, "%v", err)
	}
	return err
}

func (l *Logger) LogIfErrorP(err error, format string, args ...any) error {
	if err != nil {
		l.ErrorP(fmt.Sprintf(format, args...), "%v", err)
	}
	return err
}
