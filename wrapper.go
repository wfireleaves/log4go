// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"os"
	"strings"
)

var (
	Global Logger
)

func init() {
	Global = NewDefaultLogger(DEBUG)
}

// Wrapper for (*Logger).LoadConfiguration
func LoadConfiguration(filename string) {
	Global.LoadConfiguration(filename)
}

// Wrapper for (*Logger).AddFilter
func AddFilter(name string, lvl Level, writer LogWriter) {
	Global.AddFilter(name, lvl, writer)
}

// Wrapper for (*Logger).Close (closes and removes all logwriters)
func Close() {
	Global.Close()
}

func Crash(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(CRITICAL, strings.Repeat(" %v", len(args))[1:], args...)
	}
	panic(args)
}

// Logs the given message and crashes the program
func Crashf(format string, args ...interface{}) {
	Global.intLogf(CRITICAL, format, args...)
	Global.Close() // so that hopefully the messages get logged
	panic(fmt.Sprintf(format, args...))
}

// Compatibility with `log`
func Exit(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Compatibility with `log`
func Exitf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Compatibility with `log`
func Stderr(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Compatibility with `log`
func Stderrf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
}

// Compatibility with `log`
func Stdout(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(INFO, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Compatibility with `log`
func Stdoutf(format string, args ...interface{}) {
	Global.intLogf(INFO, format, args...)
}

// Send a log message manually
// Wrapper for (*Logger).Log
func Log(lvl Level, source, message string) {
	Global.Log(lvl, source, message)
}

// Send a formatted log message easily
// Wrapper for (*Logger).Logf
func Logf(lvl Level, format string, args ...interface{}) {
	Global.intLogf(lvl, format, args...)
}

// Send a closure log message
// Wrapper for (*Logger).Logc
func Logc(lvl Level, closure func() string) {
	Global.intLogc(lvl, closure)
}

// Utility for finest log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Finest
func Finestf(arg0 string, args ...interface{}) {
	const (
		lvl = FINEST
	)
	Global.intLogf(lvl, arg0, args...)
}

// Utility for fine log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Fine
func Finef(arg0 string, args ...interface{}) {
	const (
		lvl = FINE
	)
	Global.intLogf(lvl, arg0, args...)
}

// Utility for debug log messages
// When given a string as the first argument, this behaves like Logf but with the DEBUG log level (e.g. the first argument is interpreted as a format for the latter arguments)
// When given a closure of type func()string, this logs the string returned by the closure iff it will be logged.  The closure runs at most one time.
// When given anything else, the log message will be each of the arguments formatted with %v and separated by spaces (ala Sprint).
// Wrapper for (*Logger).Debug
func Debugf(arg0 string, args ...interface{}) {
	const (
		lvl = DEBUG
	)
	Global.intLogf(lvl, arg0, args...)
}

// Utility for trace log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Trace
func Tracef(arg0 string, args ...interface{}) {
	const (
		lvl = TRACE
	)
	Global.intLogf(lvl, arg0, args...)
}

// Utility for info log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Info
func Infof(arg0 string, args ...interface{}) {
	const (
		lvl = INFO
	)
	Global.intLogf(lvl, arg0, args...)
}

// Utility for warn log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Warn
func Warnf(arg0 string, args ...interface{}) {
	const (
		lvl = WARNING
	)
	Global.intLogf(lvl, arg0, args...)
}

// Utility for error log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Error
func Errorf(arg0 string, args ...interface{}) {
	const (
		lvl = ERROR
	)
	Global.intLogf(lvl, arg0, args...)
}

// Utility for critical log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Critical
func Criticalf(arg0 string, args ...interface{}) {
	const (
		lvl = CRITICAL
	)

	Global.intLogf(lvl, arg0, args...)
}

func Debug(message string, field ...Field) {
	const (
		lvl = DEBUG
	)
	Global.intLogJson(lvl, message, field...)
}

func Info(message string, field ...Field) {
	const (
		lvl = INFO
	)
	Global.intLogJson(lvl, message, field...)
}

func Warn(message string, field ...Field) {
	const (
		lvl = WARNING
	)
	Global.intLogJson(lvl, message, field...)
}

func Error(message string, field ...Field) {
	const (
		lvl = ERROR
	)
	Global.intLogJson(lvl, message, field...)
}
