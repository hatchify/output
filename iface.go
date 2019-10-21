package output

import (
	"context"
	"io"
	"time"
)

// ClassicOutputter represents an outputter interface from previous version of the output package.
type ClassicOutputter interface {
	Notification(format string, args ...interface{})
	Success(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
	Debug(format string, args ...interface{})
}

// Outputter represents a full outputter interface.
// It was inspired by previous ClassicOutputter interface that we have to support,
// also logrus capabilities that are added here just recently.
type Outputter interface {
	// Classic output methods
	ClassicOutputter

	// Logrus context providers

	WithField(key string, value interface{}) *Entry
	WithFields(fields Fields) *Entry
	WithError(err error) *Entry
	WithContext(ctx context.Context) *Entry
	WithTime(t time.Time) *Entry

	// Logrus formatted logging methods

	Logf(level Level, format string, args ...interface{})
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	// Logrus shortcut logging methods

	Log(level Level, args ...interface{})
	Trace(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Logln(level Level, args ...interface{})
	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	// Logrus configuration and middleware

	SetFormatter(formatter Formatter)
	SetOutput(output io.Writer)
	SetLevel(level Level)
	GetLevel() Level
	IsLevelEnabled(level Level) bool
	AddHook(hook Hook)
	ReplaceHooks(hooks LevelHooks) LevelHooks
	Exit(code int)
}

// Won't compile if StdLogger can't be realized by the outputter.
var (
	_ StdLogger = &outputter{}
	_ StdLogger = &Entry{}
)

// StdLogger is what your output-enabled library should take, that way
// it'll accept a stdlib logger (*log.Logger) and an output.Outputter. There's no standard
// interface, this is the closest we get, unfortunately.
type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}
