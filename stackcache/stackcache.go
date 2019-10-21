package stackcache

import (
	"runtime"
	"strings"
	"sync"
)

type StackCache interface {
	GetCaller() (runtime.Frame, bool)
}

// New creates a new stack cache for effectively traversing runtime callers.
//
// framesOffset param allows to have flexibility of stack trace parsing,
// by offsetting PC of the logging package entrypoint.
//
// Rule of thumb: the more "wrapped" the logging function is, the higher PC should be.
// Default for output.Outputter: 11.
// For a default outputter on package level: 12 (+1 for package-level wrappers).
// For outputter.WithFn() helper:
func New(framesOffset int) StackCache {
	return &stackCache{
		framesOffset: framesOffset,

		minimumCallerDepth: 1,
		maximumCallerDepth: 25,
	}
}

type stackCache struct {
	// qualified package name, cached at first use
	outputPackageName string

	// framesOffset sets pc offset in stack frames
	framesOffset int
	// Used for caller information initialisation
	callerInitOnce sync.Once
	// start at the bottom of the stack before the package-name cache is primed
	minimumCallerDepth int
	// limit caller depth scanning to avoid failing in weird configurations
	maximumCallerDepth int
}

// GetCaller retrieves the name of the first non-stackcache calling function.
// This function is from stackcache internals.
func (c *stackCache) GetCaller() (runtime.Frame, bool) {
	// cache this package's fully-qualified name
	c.callerInitOnce.Do(func() {
		pcs := make([]uintptr, 2)
		_ = runtime.Callers(0, pcs)
		c.outputPackageName = getPackageName(runtime.FuncForPC(pcs[1]).Name())

		// now that we have the cache, we can skip a minimum count of known-stackcache functions
		c.minimumCallerDepth = c.framesOffset
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, c.maximumCallerDepth)
	depth := runtime.Callers(c.minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of the package, we're done
		if pkg != c.outputPackageName {
			return f, true
		}
	}

	// if we got here, we failed to find the caller's context
	return runtime.Frame{}, false
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
