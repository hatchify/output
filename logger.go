package output

import (
	"fmt"
	"io"
	"sync"
)

var debugFmt = "%s:%d :: %s"

// New will return a newly created output logger
func New(wc io.WriteCloser, fn PrefixFunc) *Logger {
	var l Logger
	l.wc = wc
	l.fn = fn
	return &l
}

// Logger manages an output stream
type Logger struct {
	mux sync.RWMutex
	// Underlying write closer
	wc io.WriteCloser
	// Prefix function
	fn PrefixFunc
	// Closed state
	closed bool
}

// Print will print to the underlying writer
func (l *Logger) Print(str string) {
	l.mux.Lock()
	defer l.mux.Unlock()

	buf := make([]byte, 0, len(str)+1)
	if l.fn != nil {
		buf = append(buf, l.fn()...)
	}

	buf = append(buf, str...)
	buf = append(buf, '\n')
	l.wc.Write(buf)
}

// Printf will print a formatted message to the underlying writer
func (l *Logger) Printf(format string, values ...interface{}) {
	msg := fmt.Sprintf(format, values...)
	l.Print(msg)
}

// Notification will output a notification message
func (l *Logger) Notification(format string, values ...interface{}) {
	l.Printf(dot+format, values...)
}

// Success will output a success message
func (l *Logger) Success(format string, values ...interface{}) {
	l.Printf(greenDot+format, values...)
}

// Warning will output a warning message
func (l *Logger) Warning(format string, values ...interface{}) {
	l.Printf(yellowDot+format, values...)
}

// Error will output an error message
func (l *Logger) Error(format string, values ...interface{}) {
	l.Printf(redDot+format, values...)
}

// Debug will log a debug line
func (l *Logger) Debug(format string, values ...interface{}) {
	msg := fmt.Sprintf(format, values...)
	filename, lineNumber := getDebugVals()
	msg = fmt.Sprintf(debugFmt, filename, lineNumber, msg)
	l.Print(msg)
}

// Close will close a logger
func (l *Logger) Close() (err error) {
	l.mux.Lock()
	defer l.mux.Unlock()
	if l.closed {
		return
	}

	l.closed = true

	if err = l.wc.Close(); err != nil {
		return
	}

	return
}
