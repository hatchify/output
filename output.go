package output

import (
	"context"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	bugsnagHook "github.com/hatchify/output-bugsnag/hooks/bugsnag"
	blobHook "github.com/hatchify/output/hooks/blob"
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
		logger: &logrus.Logger{
			Out:       wc,
			Formatter: formatter,
			Hooks:     make(LevelHooks),
			Level:     DebugLevel,
			ExitFunc:  closer.Exit,
		},

		wc:       wc,
		mux:      new(sync.Mutex),
		stack:    stackcache.New(1, "github.com/hatchify/output"),
		initDone: true,
	}
	out.entry = out.logger.WithContext(context.Background())

	for _, h := range hooks {
		out.AddHook(h)
	}

	return out
}

type outputter struct {
	logger *logrus.Logger
	entry  *logrus.Entry

	mux   *sync.Mutex
	wc    io.Writer
	stack stackcache.StackCache

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
		out.logger = &logrus.Logger{
			Out:       out.wc,
			Formatter: new(TextFormatter),
			Hooks:     make(LevelHooks),
			Level:     DebugLevel,
			ExitFunc:  closer.Exit,
		}
		out.entry = out.logger.WithContext(context.Background())
		out.stack = stackcache.New(1, "github.com/hatchify/output")
		out.addDefaultHooks()
		out.mux = new(sync.Mutex)
		out.initDone = true
	})
}

// addDefaultHooks initializes default hooks and additional hooks
// based on the environment setup.
func (out *outputter) addDefaultHooks() {
	out.logger.AddHook(debugHook.NewHook(nil))

	if isTrue(os.Getenv("OUTPUT_BLOB_ENABLED")) {
		out.logger.AddHook(blobHook.NewHook(nil))
	}

	if isTrue(os.Getenv("OUTPUT_BUGSNAG_ENABLED")) {
		out.logger.AddHook(bugsnagHook.NewHook(nil))
	}
}

// Adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Error, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (out *outputter) WithField(key string, value interface{}) Outputter {
	out.initOnce()

	outCopy := out.copy()
	outCopy.entry = out.entry.WithField(key, value)

	return outCopy
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (out *outputter) WithFields(fields Fields) Outputter {
	out.initOnce()
	outCopy := out.copy()
	outCopy.entry = out.entry.WithFields(fields)

	return outCopy
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (out *outputter) WithError(err error) Outputter {
	out.initOnce()
	outCopy := out.copy()
	outCopy.entry = out.entry.WithError(err)

	return outCopy
}

// Add a context to the log entry.
func (out *outputter) WithContext(ctx context.Context) Outputter {
	out.initOnce()
	outCopy := out.copy()
	outCopy.entry = out.entry.WithContext(ctx)

	return outCopy
}

// Overrides the time of the log entry.
func (out *outputter) WithTime(t time.Time) Outputter {
	out.initOnce()
	outCopy := out.copy()
	outCopy.entry = out.entry.WithTime(t)

	return outCopy
}

func (out *outputter) Logf(level Level, format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(level, format, args...)
}

func (out *outputter) Tracef(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(TraceLevel, format, args...)
}

func (out *outputter) Debugf(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(DebugLevel, format, args...)
}

func (out *outputter) Infof(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(InfoLevel, format, args...)
}

func (out *outputter) Printf(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Printf(format, args...)
}

func (out *outputter) Warningf(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(WarnLevel, format, args...)
}

func (out *outputter) Errorf(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(ErrorLevel, format, args...)
}

func (out *outputter) Fatalf(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(FatalLevel, format, args...)
	out.logger.Exit(1)
}

func (out *outputter) Panicf(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(PanicLevel, format, args...)
}

func (out *outputter) Log(level Level, args ...interface{}) {
	out.initOnce()
	out.entry.Log(level, args...)
}

func (out *outputter) Trace(args ...interface{}) {
	out.initOnce()
	out.entry.Log(TraceLevel, args...)
}

func (out *outputter) Info(args ...interface{}) {
	out.initOnce()
	out.entry.Log(InfoLevel, args...)
}

func (out *outputter) Print(args ...interface{}) {
	out.initOnce()
	out.entry.Print(args...)
}

func (out *outputter) Fatal(args ...interface{}) {
	out.initOnce()
	out.entry.Log(FatalLevel, args...)
	out.logger.Exit(1)
}

func (out *outputter) Panic(args ...interface{}) {
	out.initOnce()
	out.entry.Log(PanicLevel, args...)
}

func (out *outputter) Logln(level Level, args ...interface{}) {
	out.initOnce()
	out.entry.Logln(level, args...)
}

func (out *outputter) Traceln(args ...interface{}) {
	out.initOnce()
	out.entry.Logln(TraceLevel, args...)
}

func (out *outputter) Debugln(args ...interface{}) {
	out.initOnce()
	out.entry.Logln(DebugLevel, args...)
}

func (out *outputter) Infoln(args ...interface{}) {
	out.initOnce()
	out.entry.Logln(InfoLevel, args...)
}

func (out *outputter) Println(args ...interface{}) {
	out.initOnce()
	out.entry.Println(args...)
}

func (out *outputter) Warningln(args ...interface{}) {
	out.initOnce()
	out.entry.Logln(WarnLevel, args...)
}

func (out *outputter) Errorln(args ...interface{}) {
	out.initOnce()
	out.entry.Logln(ErrorLevel, args...)
}

func (out *outputter) Fatalln(args ...interface{}) {
	out.initOnce()
	out.entry.Logln(FatalLevel, args...)
	out.logger.Exit(1)
}

func (out *outputter) Debug(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(DebugLevel, format, args...)
}

func (out *outputter) Notification(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(InfoLevel, format, args...)
}

func (out *outputter) Success(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(InfoLevel, format, args...)
}

func (out *outputter) Warning(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(WarnLevel, format, args...)
}

func (out *outputter) Error(format string, args ...interface{}) {
	out.initOnce()
	out.entry.Logf(ErrorLevel, format, args...)
}

func (out *outputter) Panicln(args ...interface{}) {
	out.initOnce()
	out.entry.Logln(PanicLevel, args...)
}

// SetLevel sets the logger level.
func (out *outputter) SetLevel(level Level) {
	out.initOnce()
	out.logger.SetLevel(level)
}

// GetLevel returns the logger level.
func (out *outputter) GetLevel() Level {
	out.initOnce()
	return out.logger.GetLevel()
}

// AddHook adds a hook to the logger hooks.
func (out *outputter) AddHook(hook Hook) {
	out.initOnce()
	out.logger.AddHook(hook)
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (out *outputter) IsLevelEnabled(level Level) bool {
	out.initOnce()
	return out.logger.IsLevelEnabled(level)
}

// SetFormatter sets the logger formatter.
func (out *outputter) SetFormatter(formatter Formatter) {
	out.initOnce()
	out.logger.SetFormatter(formatter)
}

// SetOutput sets the logger output.
func (out *outputter) SetOutput(output io.Writer) {
	out.initOnce()
	out.logger.SetOutput(output)
}

// ReplaceHooks replaces the logger hooks and returns the old ones
func (out *outputter) ReplaceHooks(hooks LevelHooks) LevelHooks {
	out.initOnce()
	return out.logger.ReplaceHooks(hooks)
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
	caller := out.stack.GetCaller()
	parts := strings.Split(caller.Function, "/")
	nameParts := strings.Split(parts[len(parts)-1], ".")

	return nameParts[len(nameParts)-1]
}

func isTrue(v string) bool {
	switch strings.ToLower(v) {
	case "1", "true", "y":
		return true
	}

	return false
}

// copy allows to construct an outputter copy with new entry.
func (out *outputter) copy() *outputter {
	return &outputter{
		wc:       out.wc,
		logger:   out.logger,
		stack:    out.stack,
		mux:      out.mux,
		initDone: out.initDone,
		closed:   out.closed,
	}
}
