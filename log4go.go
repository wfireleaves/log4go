// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

// Package log4go provides level-based and highly configurable logging.
//
// Enhanced Logging
//
// This is inspired by the logging functionality in Java.  Essentially, you create a Logger
// object and create output filters for it.  You can send whatever you want to the Logger,
// and it will filter that based on your settings and send it to the outputs.  This way, you
// can put as much debug code in your program as you want, and when you're done you can filter
// out the mundane messages so only the important ones show up.
//
// Utility functions are provided to make life easier. Here is some example code to get started:
//
// log := log4go.NewLogger()
// log.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())
// log.AddFilter("log",    log4go.FINE,  log4go.NewFileLogWriter("example.log", true))
// log.Info("The time is now: %s", time.LocalTime().Format("15:04:05 MST 2006/01/02"))
//
// The first two lines can be combined with the utility NewDefaultLogger:
//
// log := log4go.NewDefaultLogger(log4go.DEBUG)
// log.AddFilter("log",    log4go.FINE,  log4go.NewFileLogWriter("example.log", true))
// log.Info("The time is now: %s", time.LocalTime().Format("15:04:05 MST 2006/01/02"))
//
// Usage notes:
// - The ConsoleLogWriter does not display the source of the message to standard
//   output, but the FileLogWriter does.
// - The utility functions (Info, Debug, Warn, etc) derive their source from the
//   calling function, and this incurs extra overhead.
//
// Changes from 2.0:
// - The external interface has remained mostly stable, but a lot of the
//   internals have been changed, so if you depended on any of this or created
//   your own LogWriter, then you will probably have to update your code.  In
//   particular, Logger is now a map and ConsoleLogWriter is now a channel
//   behind-the-scenes, and the LogWrite method no longer has return values.
//
// Future work: (please let me know if you think I should work on any of these particularly)
// - Log file rotation
// - Logging configuration files ala log4j
// - Have the ability to remove filters?
// - Have GetInfoChannel, GetDebugChannel, etc return a chan string that allows
//   for another method of logging
// - Add an XML filter type
package log4go

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

// Version information
const (
	L4G_VERSION = "log4go-v3.0.1"
	L4G_MAJOR   = 3
	L4G_MINOR   = 0
	L4G_BUILD   = 1
)

/****** Constants ******/

// These are the integer logging levels used by the logger
type Level int

const (
	FINEST Level = iota
	FINE
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
)

// Logging level strings
var (
	levelStrings = [...]string{"FNST", "FINE", "DEBG", "TRAC", "INFO", "WARN", "EROR", "CRIT"}
)

func (l Level) String() string {
	if l < 0 || int(l) > len(levelStrings) {
		return "UNKNOWN"
	}
	return levelStrings[int(l)]
}

/****** Variables ******/
var (
	// LogBufferLength specifies how many log messages a particular log4go
	// logger can buffer at a time before writing them.
	LogBufferLength = 32
	LogRecordPool   = sync.Pool{New: func() interface{} {
		return newLogRecord()
	}}
)

/****** LogRecord ******/

// A LogRecord contains all of the pertinent information for each message
type LogRecord struct {
	Level   Level     // The log level
	Created time.Time // The time at which the log message was created (nanoseconds)
	Source  string    // The message source
	Message string    // The log message
	Json    bool      // The log type (true: Format, false: json)
	Fields  []Field   // The json log field
}

func newLogRecord() *LogRecord {
	return &LogRecord{}
}

func GetLogRecord(lv Level, score, message string, json bool, field []Field) *LogRecord {
	rec := LogRecordPool.Get().(*LogRecord)
	rec.Level = lv
	rec.Created = time.Now()
	rec.Source = score
	rec.Message = message
	rec.Json = json
	rec.Fields = field
	return rec
}

func PutLogRecord(rec *LogRecord) {
	rec.Level = FINEST
	rec.Created = time.Time{}
	rec.Source = ""
	rec.Message = ""
	rec.Json = false
	rec.Fields = nil
	LogRecordPool.Put(rec)
}

/****** LogWriter ******/

// This is an interface for anything that should be able to write logs
type LogWriter interface {
	// This will be called to log a LogRecord message.
	LogWrite(rec *LogRecord)

	// This should clean up anything lingering about the LogWriter, as it is called before
	// the LogWriter is removed.  LogWrite should not be called after Close.
	Close()
}

/****** Logger ******/

// A Filter represents the log level below which no log records are written to
// the associated LogWriter.
type Filter struct {
	Level Level
	LogWriter
}

// A Logger represents a collection of Filters through which log messages are
// written.
type Logger map[string]*Filter

// Create a new logger.
//
// DEPRECATED: Use make(Logger) instead.
func NewLogger() Logger {
	os.Stderr.WriteString("warning: use of deprecated NewLogger\n")
	return make(Logger)
}

// Create a new logger with a "stdout" filter configured to send log messages at
// or above lvl to standard output.
//
// DEPRECATED: use NewDefaultLogger instead.
func NewConsoleLogger(lvl Level) Logger {
	os.Stderr.WriteString("warning: use of deprecated NewConsoleLogger\n")
	return Logger{
		"stdout": &Filter{lvl, NewConsoleLogWriter()},
	}
}

// Create a new logger with a "stdout" filter configured to send log messages at
// or above lvl to standard output.
func NewDefaultLogger(lvl Level) Logger {
	return Logger{
		"stdout": &Filter{lvl, NewConsoleLogWriter()},
	}
}

// Closes all log writers in preparation for exiting the program or a
// reconfiguration of logging.  Calling this is not really imperative, unless
// you want to guarantee that all log messages are written.  Close removes
// all filters (and thus all LogWriters) from the logger.
func (log Logger) Close() {
	// Close all open loggers
	for name, filt := range log {
		filt.Close()
		delete(log, name)
	}
}

// Add a new LogWriter to the Logger which will only log messages at lvl or
// higher.  This function should not be called from multiple goroutines.
// Returns the logger for chaining.
func (log Logger) AddFilter(name string, lvl Level, writer LogWriter) Logger {
	log[name] = &Filter{lvl, writer}
	return log
}

/******* Logging *******/
// Send a formatted log message internally
func (log Logger) intLogf(lvl Level, format string, args ...interface{}) {
	skip := true

	// Determine if any logging will be done
	for _, filt := range log {
		if lvl >= filt.Level {
			skip = false
			break
		}
	}
	if skip {
		return
	}

	// Determine caller func
	_, fileName, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", fileName, lineno)
	}

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	// Make the log record
	rec := &LogRecord{
		Level:   lvl,
		Created: time.Now(),
		Source:  src,
		Message: msg,
	}

	// Dispatch the logs
	for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		filt.LogWrite(rec)
	}
}

func (log Logger) intLogJson(lvl Level, message string, filed ...Field) {
	skip := true

	// Determine if any logging will be done
	for _, filt := range log {
		if lvl >= filt.Level {
			skip = false
			break
		}
	}
	if skip {
		return
	}

	// Determine caller func
	_, fileName, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", fileName, lineno)
	}
	now := time.Now()
	// Dispatch the logs
	for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		// Make the log record
		rec := &LogRecord{
			Level:   lvl,
			Created: now,
			Source:  src,
			Message: message,
			Json:    true,
			Fields:  filed,
		}
		filt.LogWrite(rec)
	}
}

// Send a closure log message internally
func (log Logger) intLogc(lvl Level, closure func() string) {
	skip := true

	// Determine if any logging will be done
	for _, filt := range log {
		if lvl >= filt.Level {
			skip = false
			break
		}
	}
	if skip {
		return
	}

	// Determine caller func
	_, fileName, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", fileName, lineno)
	}

	// Make the log record
	rec := &LogRecord{
		Level:   lvl,
		Created: time.Now(),
		Source:  src,
		Message: closure(),
	}

	// Dispatch the logs
	for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		filt.LogWrite(rec)
	}
}

// Send a log message with manual level, source, and message.
func (log Logger) Log(lvl Level, source, message string) {
	skip := true

	// Determine if any logging will be done
	for _, filt := range log {
		if lvl >= filt.Level {
			skip = false
			break
		}
	}
	if skip {
		return
	}

	// Make the log record
	rec := &LogRecord{
		Level:   lvl,
		Created: time.Now(),
		Source:  source,
		Message: message,
	}

	// Dispatch the logs
	for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		filt.LogWrite(rec)
	}
}

// Logf logs a formatted log message at the given log level, using the caller as
// its source.
func (log Logger) Logf(lvl Level, format string, args ...interface{}) {
	log.intLogf(lvl, format, args...)
}

// Logc logs a string returned by the closure at the given log level, using the caller as
// its source.  If no log message would be written, the closure is never called.
func (log Logger) Logc(lvl Level, closure func() string) {
	log.intLogc(lvl, closure)
}

// Finestf logs a message at the finest log level.
// See Debug for an explanation of the arguments.
func (log Logger) Finestf(arg0 string, args ...interface{}) {
	const (
		lvl = FINEST
	)
	log.intLogf(lvl, arg0, args...)
}

// Finef logs a message at the fine log level.
// See Debug for an explanation of the arguments.
func (log Logger) Finef(arg0 string, args ...interface{}) {
	const (
		lvl = FINE
	)
	log.intLogf(lvl, arg0, args...)
}

// Debugf is a utility method for debug log messages.
// The behavior of Debug depends on the first argument:
// - arg0 is a string
//   When given a string as the first argument, this behaves like Logf but with
//   the DEBUG log level: the first argument is interpreted as a format for the
//   latter arguments.
// - arg0 is a func()string
//   When given a closure of type func()string, this logs the string returned by
//   the closure iff it will be logged.  The closure runs at most one time.
// - arg0 is interface{}
//   When given anything else, the log message will be each of the arguments
//   formatted with %v and separated by spaces (ala Sprint).
func (log Logger) Debugf(arg0 string, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	log.intLogf(lvl, arg0, args...)
}

// Tracef logs a message at the trace log level.
// See Debug for an explanation of the arguments.
func (log Logger) Tracef(arg0 string, args ...interface{}) {
	const (
		lvl = TRACE
	)
	log.intLogf(lvl, arg0, args...)
}

// Infof logs a message at the info log level.
// See Debug for an explanation of the arguments.
func (log Logger) Infof(arg0 string, args ...interface{}) {
	const (
		lvl = INFO
	)
	log.intLogf(lvl, arg0, args...)
}

// Warnf logs a message at the warning log level and returns the formatted error.
// At the warning level and higher, there is no performance benefit if the
// message is not actually logged, because all formats are processed and all
// closures are executed to format the error message.
// See Debug for further explanation of the arguments.
func (log Logger) Warnf(arg0 string, args ...interface{}) {
	const (
		lvl = WARNING
	)
	log.intLogf(lvl, fmt.Sprintf(arg0, args...))
}

// Errorf logs a message at the error log level and returns the formatted error,
// See Warn for an explanation of the performance and Debug for an explanation
// of the parameters.
func (log Logger) Errorf(arg0 string, args ...interface{}) {
	const (
		lvl = ERROR
	)
	log.intLogf(lvl, fmt.Sprintf(arg0, args...))
}

// Criticalf logs a message at the critical log level and returns the formatted error,
// See Warn for an explanation of the performance and Debug for an explanation
// of the parameters.
func (log Logger) Criticalf(arg0 string, args ...interface{}) {
	const (
		lvl = CRITICAL
	)
	log.intLogf(lvl, fmt.Sprintf(arg0, args...))
}

func (log Logger) Debug(message string, field ...Field) {
	const (
		lvl = DEBUG
	)
	log.intLogJson(lvl, message, field...)
}

func (log Logger) Info(message string, field ...Field) {
	const (
		lvl = INFO
	)
	log.intLogJson(lvl, message, field...)
}

func (log Logger) Warn(message string, field ...Field) {
	const (
		lvl = WARNING
	)
	log.intLogJson(lvl, message, field...)
}

func (log Logger) Error(message string, field ...Field) {
	const (
		lvl = ERROR
	)
	log.intLogJson(lvl, message, field...)
}
