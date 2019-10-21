package output

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	debugHook "github.com/hatchify/output/hooks/debug"
	"github.com/hatchify/output/stackcache"
	"github.com/sirupsen/logrus"
	"github.com/xlab/closer"
)

// NewOutputter constructs a new outputter.
func NewOutputter(wc io.Writer, formatter Formatter, hooks ...Hook) Outputter {
	if formatter == nil {
		formatter = new(TextFormatter)
	}
	out := &outputter{
		Logger: &logrus.Logger{
			Out:       wc,
			Formatter: formatter,
			Hooks:     make(LevelHooks),
			Level:     DebugLevel,
			ExitFunc:  closer.Exit,
		},

		wc:       wc,
		mux:      new(sync.Mutex),
		stack:    stackcache.New(3),
		initDone: true,
	}
	for _, h := range hooks {
		out.AddHook(h)
	}
	return out
}

type outputter struct {
	*logrus.Logger

	mux         *sync.Mutex
	wc          io.Writer
	stack       stackcache.StackCache
	stackOffset int

	init     sync.Once
	initDone bool
	closed   bool
}

func (out *outputter) initOnce() {
	out.init.Do(func() {
		if out.initDone {
			// bail out if init already done (if New contstructor has been used).
			return
		}
		if out.wc == nil {
			out.wc = os.Stderr
		}
		// otherwise init output with conservative defaults
		out.Logger = &logrus.Logger{
			Out:       out.wc,
			Formatter: new(TextFormatter),
			Hooks:     make(LevelHooks),
			Level:     DebugLevel,
			ExitFunc:  closer.Exit,
		}
		if out.stackOffset == 0 {
			out.stackOffset = 4
		}
		out.stack = stackcache.New(out.stackOffset)
		out.Logger.AddHook(debugHook.NewHook(&debugHook.HookOptions{
			FramesOffset: 12,
		}))
		out.mux = new(sync.Mutex)
		out.initDone = true
	})
}

// Adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Error, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (out *outputter) WithField(key string, value interface{}) *Entry {
	out.initOnce()
	return out.Logger.WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (out *outputter) WithFields(fields Fields) *Entry {
	out.initOnce()
	return out.Logger.WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (out *outputter) WithError(err error) *Entry {
	out.initOnce()
	return out.Logger.WithError(err)
}

// Add a context to the log entry.
func (out *outputter) WithContext(ctx context.Context) *Entry {
	out.initOnce()
	return out.Logger.WithContext(ctx)
}

// Overrides the time of the log entry.
func (out *outputter) WithTime(t time.Time) *Entry {
	out.initOnce()
	return out.Logger.WithTime(t)
}

func (out *outputter) Logf(level Level, format string, args ...interface{}) {
	out.initOnce()
	out.Logger.Logf(level, format, args...)
}

func (out *outputter) Tracef(format string, args ...interface{}) {
	out.Logf(TraceLevel, format, args...)
}

func (out *outputter) Debugf(format string, args ...interface{}) {
	out.initOnce()
	out.Logf(DebugLevel, format, args...)
}

func (out *outputter) Infof(format string, args ...interface{}) {
	out.Logf(InfoLevel, format, args...)
}

func (out *outputter) Printf(format string, args ...interface{}) {
	out.initOnce()
	out.Logger.Printf(format, args...)
}

func (out *outputter) Warningf(format string, args ...interface{}) {
	out.initOnce()
	out.Logf(WarnLevel, format, args...)
}

func (out *outputter) Fatalf(format string, args ...interface{}) {
	out.Logf(FatalLevel, format, args...)
	out.Exit(1)
}

func (out *outputter) Panicf(format string, args ...interface{}) {
	out.Logf(PanicLevel, format, args...)
}

func (out *outputter) Log(level Level, args ...interface{}) {
	out.initOnce()
	out.Logger.Log(level, args...)
}

func (out *outputter) Trace(args ...interface{}) {
	out.Log(TraceLevel, args...)
}

func (out *outputter) Info(args ...interface{}) {
	out.Log(InfoLevel, args...)
}

func (out *outputter) Print(args ...interface{}) {
	out.initOnce()
	out.Logger.Print(args...)
}

func (out *outputter) Fatal(args ...interface{}) {
	out.Log(FatalLevel, args...)
	out.Exit(1)
}

func (out *outputter) Panic(args ...interface{}) {
	out.Log(PanicLevel, args...)
}

func (out *outputter) Logln(level Level, args ...interface{}) {
	out.initOnce()
	out.Logger.Logln(level, args...)
}

func (out *outputter) Traceln(args ...interface{}) {
	out.Logln(TraceLevel, args...)
}

func (out *outputter) Debugln(args ...interface{}) {
	out.initOnce()
	out.Logln(DebugLevel, args...)
}

func (out *outputter) Infoln(args ...interface{}) {
	out.Logln(InfoLevel, args...)
}

func (out *outputter) Println(args ...interface{}) {
	out.initOnce()
	out.Logger.Println(args...)
}

func (out *outputter) Warningln(args ...interface{}) {
	out.Logln(WarnLevel, args...)
}

func (out *outputter) Errorln(args ...interface{}) {
	out.Logln(ErrorLevel, args...)
}

func (out *outputter) Fatalln(args ...interface{}) {
	out.Logln(FatalLevel, args...)
	out.Exit(1)
}

func (out *outputter) Debug(format string, args ...interface{}) {
	out.initOnce()
	out.Logf(DebugLevel, format, args...)
}

func (out *outputter) Notification(format string, args ...interface{}) {
	out.initOnce()
	out.Logf(InfoLevel, format, args...)
}

func (out *outputter) Success(format string, args ...interface{}) {
	out.initOnce()
	out.Logf(SuccessLevel, format, args...)
}

func (out *outputter) Warning(format string, args ...interface{}) {
	out.initOnce()
	out.Logf(WarnLevel, format, args...)
}

func (out *outputter) Error(format string, args ...interface{}) {
	out.initOnce()
	out.Logf(ErrorLevel, format, args...)
}

func (out *outputter) Exit(code int) {
	out.initOnce()
	out.Logger.Exit(code)
}

func (out *outputter) Panicln(args ...interface{}) {
	out.Logln(PanicLevel, args...)
}

// SetLevel sets the logger level.
func (out *outputter) SetLevel(level Level) {
	out.initOnce()
	out.Logger.SetLevel(level)
}

// GetLevel returns the logger level.
func (out *outputter) GetLevel() Level {
	out.initOnce()
	return out.Logger.GetLevel()
}

// AddHook adds a hook to the logger hooks.
func (out *outputter) AddHook(hook Hook) {
	out.initOnce()
	out.Logger.AddHook(hook)
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (out *outputter) IsLevelEnabled(level Level) bool {
	out.initOnce()
	return out.Logger.IsLevelEnabled(level)
}

// SetFormatter sets the logger formatter.
func (out *outputter) SetFormatter(formatter Formatter) {
	out.initOnce()
	out.Logger.SetFormatter(formatter)
}

// SetOutput sets the logger output.
func (out *outputter) SetOutput(output io.Writer) {
	out.initOnce()
	out.Logger.SetOutput(output)
}

// ReplaceHooks replaces the logger hooks and returns the old ones
func (out *outputter) ReplaceHooks(hooks LevelHooks) LevelHooks {
	out.initOnce()
	return out.Logger.ReplaceHooks(hooks)
}

// Close effectively closes output, closing the underlying writer
// if it implements io.WriteCloser.
func (out *outputter) Close() (err error) {
	// bail out if already closed
	out.mux.Lock()
	defer out.mux.Unlock()
	if out.closed {
		return
	}
	out.closed = true

	// try to close only WriteClosers
	if outCloser, ok := out.wc.(io.WriteCloser); ok {
		return outCloser.Close()
	}
	return
}

// CallerName returns caller function name.
func (out *outputter) CallerName() string {
	out.initOnce()
	caller, ok := out.stack.GetCaller()
	if !ok {
		return ""
	}
	parts := strings.Split(caller.Function, "/")
	nameParts := strings.Split(parts[len(parts)-1], ".")
	return nameParts[len(nameParts)-1]
}
