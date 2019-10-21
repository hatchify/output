package debug

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// HookOptions allows to set additional Hook options.
type HookOptions struct {
	// Levels enables this hook for all listed levels.
	Levels []logrus.Level
	// FramesOffset allows to have flexibility of stack trace parsing,
	// by offsetting PC of the logging package entrypoint.
	//
	// Rule of thumb: the more "wrapped" the logging function is, the higher PC should be.
	// Default for output.Outputter: 11.
	// For a default outputter on package level: 12 (+1 for package-level wrappers).
	FramesOffset int
	// PathSegmentsLimit allows to trim amount of source code file path segments.
	// Untrimmed: /Users/xlab/Documents/dev/go/src/github.com/hatchify/output/default_test.go
	// Trimmed (3): hatchify/output/default_test.go
	PathSegmentsLimit int
}

func checkHookOptions(opt *HookOptions) *HookOptions {
	if opt == nil {
		opt = &HookOptions{}
	}
	if len(opt.Levels) == 0 {
		opt.Levels = []logrus.Level{
			logrus.DebugLevel,
			logrus.TraceLevel,
		}
	}
	if opt.FramesOffset == 0 {
		opt.FramesOffset = 11
	}
	if opt.PathSegmentsLimit == 0 {
		opt.PathSegmentsLimit = 3
	}
	return opt
}

// NewHook initializes a new logrus.Hook using provided params and options.
func NewHook(opt *HookOptions) logrus.Hook {
	opt = checkHookOptions(opt)
	return &hook{
		opt: opt,
		// start at the bottom of the stack before the package-name cache is primed
		minimumCallerDepth: 1,
		// limit caller depth scanning to avoid failing in weird configurations
		maximumCallerDepth: 25,
	}
}

type hook struct {
	opt *HookOptions

	// qualified package name, cached at first use
	outputPackageName string
	// Used for caller information initialisation
	callerInitOnce sync.Once

	minimumCallerDepth int
	maximumCallerDepth int
}

func (h *hook) Levels() []logrus.Level {
	return h.opt.Levels
}

func (h *hook) Fire(e *logrus.Entry) error {
	caller, ok := h.getCaller()
	if !ok {
		// no caller info
		return nil
	}
	if len(caller.Function) > 0 {
		parts := strings.Split(caller.Function, "/")
		nameParts := strings.Split(parts[len(parts)-1], ".")
		e.Data["fn"] = nameParts[len(nameParts)-1]
	}
	callerFile := limitPath(caller.File, h.opt.PathSegmentsLimit)
	e.Data["src"] = fmt.Sprintf("%s:%d", callerFile, caller.Line)
	return nil
}

func limitPath(path string, n int) string {
	if n <= 0 {
		return path
	}
	pathParts := strings.Split(path, string(filepath.Separator))
	if len(pathParts) > n {
		pathParts = pathParts[len(pathParts)-n:]
	}
	return filepath.Join(pathParts...)
}

// getPackageName reduces a fully qualified function name to the package name
// This function is from logrus internals.
func getPackageName(path string) string {
	for {
		lastPeriod := strings.LastIndex(path, ".")
		lastSlash := strings.LastIndex(path, "/")
		if lastPeriod > lastSlash {
			path = path[:lastPeriod]
		} else {
			break
		}
	}
	return path
}

// getCaller retrieves the name of the first non-logrus calling function.
// This function is from logrus internals.
func (h *hook) getCaller() (runtime.Frame, bool) {
	// cache this package's fully-qualified name
	h.callerInitOnce.Do(func() {
		pcs := make([]uintptr, 2)
		_ = runtime.Callers(0, pcs)
		h.outputPackageName = getPackageName(runtime.FuncForPC(pcs[1]).Name())

		// now that we have the cache, we can skip a minimum count of known-logrus functions
		h.minimumCallerDepth = h.opt.FramesOffset
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, h.maximumCallerDepth)
	depth := runtime.Callers(h.minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of the package, we're done
		if pkg != h.outputPackageName {
			return f, true
		}
	}

	// if we got here, we failed to find the caller's context
	return runtime.Frame{}, false
}
