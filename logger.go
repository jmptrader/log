package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

///////////////////////////
//                       //
// Logger implementation //
//                       //
///////////////////////////

// logger is the default Logger implementation.
type logger struct {
	Name      string
	Level     Level
	Handler   Handler
	calldepth int
}

// NewLogger returns a new Logger implementation. Do not forget to close it at exit.
func NewLogger(name string) Logger {
	return &logger{
		Name:    name,
		Level:   DefaultLevel,
		Handler: DefaultHandler,
	}
}

func (l *logger) SetLevel(level Level) {
	l.Level = level
}

func (l *logger) SetHandler(b Handler) {
	l.Handler = b
}

func (l *logger) SetCallDepth(n int) {
	l.calldepth = n
}

func (l *logger) log(level Level, format *string, args ...interface{}) {
	if level > l.Level {
		return
	}

	_, file, line, ok := runtime.Caller(l.calldepth + 2)
	if !ok {
		file = "???"
		line = 0
	}

	rec := &Record{
		format:      format,
		args:        args,
		LoggerName:  l.Name,
		Level:       level,
		Time:        time.Now(),
		Filename:    file,
		Line:        line,
		ProcessName: procName(),
		ProcessID:   os.Getpid(),
	}

	if format != nil {
		rec.Message = fmt.Sprintf(*format, args...)
	} else {
		rec.Message = fmt.Sprint(args...)
	}

	// Add missing newline at the end.
	if !strings.HasSuffix(rec.Message, "\n") {
		rec.Message += "\n"
	}

	l.Handler.Handle(rec)
}

// procName returns the name of the current process.
func procName() string { return filepath.Base(os.Args[0]) }

// Fatal is equivalent to Critical() followed by a call to os.Exit(1).
func (l *logger) Fatal(args ...interface{}) {
	l.Critical(args...)
	os.Exit(1)
}

// Fatalf is equivalent to Criticalf() followed by a call to os.Exit(1).
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.Criticalf(format, args...)
	os.Exit(1)
}

// Panic is equivalent to Critical() followed by a call to panic().
func (l *logger) Panic(args ...interface{}) {
	l.Critical(args...)
	panic(fmt.Sprint(args...))
}

// Panicf is equivalent to Criticalf() followed by a call to panic().
func (l *logger) Panicf(format string, args ...interface{}) {
	l.Criticalf(format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Critical sends a critical level log message to the handler.
func (l *logger) Critical(args ...interface{}) { l.log(CRITICAL, nil, args...) }

// Criticalf sends a critical level log message to the handler. Arguments are handled in the manner of fmt.Printf.
func (l *logger) Criticalf(format string, args ...interface{}) { l.log(CRITICAL, &format, args...) }

// Error sends a error level log message to the handler.
func (l *logger) Error(args ...interface{}) { l.log(ERROR, nil, args...) }

// Errorf sends a error level log message to the handler. Arguments are handled in the manner of fmt.Printf.
func (l *logger) Errorf(format string, args ...interface{}) { l.log(ERROR, &format, args...) }

// Warning sends a warning level log message to the handler.
func (l *logger) Warning(args ...interface{}) { l.log(WARNING, nil, args...) }

// Warningf sends a warning level log message to the handler. Arguments are handled in the manner of fmt.Printf.
func (l *logger) Warningf(format string, args ...interface{}) { l.log(WARNING, &format, args...) }

// Notice sends a notice level log message to the handler.
func (l *logger) Notice(args ...interface{}) { l.log(NOTICE, nil, args...) }

// Noticef sends a notice level log message to the handler. Arguments are handled in the manner of fmt.Printf.
func (l *logger) Noticef(format string, args ...interface{}) { l.log(NOTICE, &format, args...) }

// Info sends a info level log message to the handler.
func (l *logger) Info(args ...interface{}) { l.log(INFO, nil, args...) }

// Infof sends a info level log message to the handler. Arguments are handled in the manner of fmt.Printf.
func (l *logger) Infof(format string, args ...interface{}) { l.log(INFO, &format, args...) }

// Debug sends a debug level log message to the handler.
func (l *logger) Debug(args ...interface{}) { l.log(DEBUG, nil, args...) }

// Debugf sends a debug level log message to the handler. Arguments are handled in the manner of fmt.Printf.
func (l *logger) Debugf(format string, args ...interface{}) { l.log(DEBUG, &format, args...) }
