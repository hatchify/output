package output

import (
	"io"
	"sync"
)

// NewWrapper will return a wrapper over default logger. This is for compatibility with
// "classic" output package. Please use Outputter and its WithFields method for any new code.
func NewWrapper(prefix string) *Wrapper {
	return NewWrapperWithOutputter(defaultOut, prefix)
}

// NewWrapperWithOutputter will return a new wrapper with a custom Logger.
// This is for compatibility with "classic" output package.
// Please use Outputter and its WithFields method for any new code.
func NewWrapperWithOutputter(o Outputter, prefix string) *Wrapper {
	w := &Wrapper{
		Entry:  o.WithField("prefix", prefix),
		prefix: prefix,
	}
	return w
}

// Wrapper will wrap an output Entry with prefix.
// You must avoid using Wrapper in any new code, use WithFields instead.
type Wrapper struct {
	*Entry
	prefix string
}

// PrefixFunc is a classic output logger prefix function.
// Use WithFields for any new code.
type PrefixFunc func() string

// New will return a classic output logger.
func New(wc io.WriteCloser, prefixFn PrefixFunc) *Logger {
	outForClassic := NewOutputter(wc, new(TextFormatter))
	l := &Logger{
		Entry:    outForClassic.WithField("logger", "classic"),
		wc:       wc,
		mux:      new(sync.Mutex),
		prefixFn: prefixFn,
	}
	return l
}

// Logger is an instance of ClassicOutputter, it manages an output stream.
// You should migrate from Logger to output.Outputter iface for any new code.
type Logger struct {
	*Entry

	wc       io.WriteCloser
	mux      *sync.Mutex
	prefixFn PrefixFunc
	closed   bool
}

func (out *Logger) entryWithPrefix() *Entry {
	if out.prefixFn != nil {
		return out.Entry.WithField("prefix", out.prefixFn())
	}
	return out.Entry
}

func (out *Logger) Debug(format string, args ...interface{}) {
	out.entryWithPrefix().Logf(DebugLevel, format, args...)
}

func (out *Logger) Notification(format string, args ...interface{}) {
	out.entryWithPrefix().Logf(InfoLevel, format, args...)
}

func (out *Logger) Print(msg string) {
	out.entryWithPrefix().Log(InfoLevel, msg)
}

func (out *Logger) Printf(format string, args ...interface{}) {
	out.entryWithPrefix().Logf(InfoLevel, format, args...)
}

func (out *Logger) Success(format string, args ...interface{}) {
	out.entryWithPrefix().Logf(SuccessLevel, format, args...)
}

func (out *Logger) Warning(format string, args ...interface{}) {
	out.entryWithPrefix().Logf(WarnLevel, format, args...)
}

func (out *Logger) Error(format string, args ...interface{}) {
	out.entryWithPrefix().Logf(ErrorLevel, format, args...)
}

func (out *Logger) Close() (err error) {
	out.mux.Lock()
	defer out.mux.Unlock()
	if out.closed {
		return
	}
	out.closed = true
	return out.wc.Close()
}
