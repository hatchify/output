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

// PrefixFunc is used to set a prefix to a log line
type PrefixFunc func() string

// SetRFC3339Timestamp is a helper function for setting RFC3339 timestamp
func SetRFC3339Timestamp() string {
	return time.Now().Format(time.RFC3339) + " :: "
}
