package output

import (
	"context"
	"time"
)

var (
	//nolint:gochecknoglobals
	defaultOut           = &outputter{}
	_          Outputter = defaultOut
)

// CLASSIC LOGGER METHODS

// Print will print to the underlying writer.
func Print(str string) {
	defaultOut.Print(str)
}

// Printf will print a formatted message to the underlying writer.
func Printf(format string, args ...interface{}) {
	defaultOut.Printf(format, args...)
}

// Notification will output a notification message.
func Notification(format string, args ...interface{}) {
	defaultOut.Notification(format, args...)
}

// Success will output a success message.
func Success(format string, args ...interface{}) {
	defaultOut.Success(format, args...)
}

// Warning will output a warning message.
func Warning(format string, args ...interface{}) {
	defaultOut.Warning(format, args...)
}

// Error will output an error message.
func Error(format string, args ...interface{}) {
	defaultOut.Error(format, args...)
}

// Debug will log a debug line.
func Debug(format string, args ...interface{}) {
	defaultOut.Debug(format, args...)
}

// OUTPUTTER METHODS
//
// Part A: Context providers

func WithField(key string, value interface{}) Outputter {
	return defaultOut.WithField(key, value)
}

func WithFields(fields Fields) Outputter {
	return defaultOut.WithFields(fields)
}

func WithError(err error) Outputter {
	return defaultOut.WithError(err)
}

func WithContext(ctx context.Context) Outputter {
	return defaultOut.WithContext(ctx)
}

func WithTime(t time.Time) Outputter {
	return defaultOut.WithTime(t)
}

// Part B: Formatted logging methods

func Logf(level Level, format string, args ...interface{}) {
	defaultOut.Logf(level, format, args...)
}

func Tracef(format string, args ...interface{}) {
	defaultOut.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	defaultOut.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultOut.Infof(format, args...)
}

func Warningf(format string, args ...interface{}) {
	defaultOut.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultOut.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultOut.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	defaultOut.Panicf(format, args...)
}

// Part C: Shortcut logging methods

func Log(level Level, args ...interface{}) {
	defaultOut.Log(level, args...)
}

func Trace(args ...interface{}) {
	defaultOut.Trace(args...)
}

func Info(args ...interface{}) {
	defaultOut.Info(args...)
}

func Fatal(args ...interface{}) {
	defaultOut.Fatal(args...)
}

func Panic(args ...interface{}) {
	defaultOut.Panic(args...)
}

func Logln(level Level, args ...interface{}) {
	defaultOut.Logln(level, args...)
}

func Traceln(args ...interface{}) {
	defaultOut.Traceln(args...)
}

func Debugln(args ...interface{}) {
	defaultOut.Debugln(args...)
}

func Infoln(args ...interface{}) {
	defaultOut.Infoln(args...)
}

func Println(args ...interface{}) {
	defaultOut.Println(args...)
}

func Warningln(args ...interface{}) {
	defaultOut.Warningln(args...)
}

func Errorln(args ...interface{}) {
	defaultOut.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	defaultOut.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	defaultOut.Panicln(args...)
}

func FnName() string {
	return defaultOut.CallerName()
}
