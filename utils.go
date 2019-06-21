package output

import (
	"path"
	"runtime"
	"time"
)

func getDebugVals() (filename string, lineNumber int) {
	_, filename, lineNumber, _ = runtime.Caller(3)
	filename = path.Clean(filename)
	return
}

// SetRFC3339Timestamp is a helper function for setting RFC3339 timestamp
func SetRFC3339Timestamp() string {
	return time.Now().Format(time.RFC3339) + " :: "
}

// PrefixFunc is used to set a prefix to a log line
type PrefixFunc func() string

// Outputter represents an outputter interface
type Outputter interface {
	Print(msg string)
	Printf(format string, values ...interface{})
	Notification(format string, values ...interface{})
	Success(format string, values ...interface{})
	Warning(format string, values ...interface{})
	Error(format string, values ...interface{})
	Debug(format string, values ...interface{})
}
