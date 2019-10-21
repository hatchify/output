package debug

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hatchify/output/stackcache"
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
		opt:   opt,
		stack: stackcache.New(opt.FramesOffset),
	}
}

type hook struct {
	opt   *HookOptions
	stack stackcache.StackCache
}

func (h *hook) Levels() []logrus.Level {
	return h.opt.Levels
}

func (h *hook) Fire(e *logrus.Entry) error {
	caller, ok := h.stack.GetCaller()
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
