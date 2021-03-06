package output

import (
	"io"
	"os"
	"sync"

	bugsnagHook "github.com/hatchify/output-bugsnag/hooks/bugsnag"
	blobHook "github.com/hatchify/output/hooks/blob"
	debugHook "github.com/hatchify/output/hooks/debug"
)

// NewWrapper will return a wrapper over default logger. This is for compatibility with
// "classic" output package. Please use Outputter and its WithFields method for any new code.
func NewWrapper(prefix string) *Wrapper {
	return NewWrapperWithOutputter(defaultOut, prefix)
}

// NewWrapperWithOutputter will return a new wrapper with a custom Logger.
// This is for compatibility with "classic" output package.
// Please use Outputter and its WithFields method for any new code.
func NewWrapperWithOutputter(o Outputter, prefix string) (w *Wrapper) {
	w = new(Wrapper)
	w.Outputter = o.WithField("prefix", prefix)
	w.prefix = prefix

	return
}

// Wrapper will wrap an output Entry with prefix.
// You must avoid using Wrapper in any new code, use WithFields instead.
type Wrapper struct {
	Outputter
	prefix string
}

// PrefixFunc is a classic output logger prefix function.
// Use WithFields for any new code.
type PrefixFunc func() string

// New will return a classic output logger.
func New(wc io.WriteCloser, prefixFn PrefixFunc) (l *Logger) {
	formatter := new(TextFormatter)
	formatter.ForceColors = true
	formatter.DisableTimestamp = true

	outForClassic := NewOutputter(wc, formatter)
	l = new(Logger)
	l.out = outForClassic.WithField("logger", "classic").(*outputter)
	l.wc = wc
	l.mux = new(sync.Mutex)
	l.prefixFn = prefixFn

	l.addDefaultHooks()

	return
}

// Logger is an instance of ClassicOutputter, it manages an output stream.
// You should migrate from Logger to output.Outputter iface for any new code.
type Logger struct {
	out      *outputter
	wc       io.WriteCloser
	mux      *sync.Mutex
	prefixFn PrefixFunc
	closed   bool
}

func (l *Logger) outWithPrefix() Outputter {
	if l.prefixFn != nil {
		return l.out.WithField("prefix", l.prefixFn())
	}

	return l.out
}

// addDefaultHooks initializes default hooks and additional hooks
// based on the environment setup.
func (l *Logger) addDefaultHooks() {
	l.out.AddHook(debugHook.NewHook(nil))

	if isTrue(os.Getenv("OUTPUT_BLOB_ENABLED")) {
		hook, _ := blobHook.NewHook(nil)
		l.out.AddHook(hook)
	}

	if isTrue(os.Getenv("OUTPUT_BUGSNAG_ENABLED")) {
		l.out.AddHook(bugsnagHook.NewHook(nil))
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.outWithPrefix().Logf(DebugLevel, format, args...)
}

func (l *Logger) Notification(format string, args ...interface{}) {
	l.outWithPrefix().Logf(InfoLevel, format, args...)
}

func (l *Logger) Print(msg string) {
	l.outWithPrefix().Log(InfoLevel, msg)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.outWithPrefix().Logf(InfoLevel, format, args...)
}

func (l *Logger) Success(format string, args ...interface{}) {
	l.outWithPrefix().Logf(InfoLevel, format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.outWithPrefix().Logf(WarnLevel, format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.outWithPrefix().Logf(ErrorLevel, format, args...)
}

func (l *Logger) Close() (err error) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.closed {
		return
	}

	l.closed = true

	return l.wc.Close()
}
