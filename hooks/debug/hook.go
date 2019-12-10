package debug

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hatchify/output/stackcache"
	"github.com/sirupsen/logrus"
)

// HookOptions allows to set additional Hook options.
type HookOptions struct {
	// AppVersion specifies version of the app currently running.
	AppVersion string
	// Levels enables this hook for all listed levels.
	Levels []logrus.Level
	// PathSegmentsLimit allows to trim amount of source code file path segments.
	// Untrimmed: /Users/xlab/Documents/dev/go/src/github.com/hatchify/output/default_test.go
	// Trimmed (3): hatchify/output/default_test.go
	PathSegmentsLimit int
}

func checkHookOptions(opt *HookOptions) *HookOptions {
	if opt == nil {
		opt = &HookOptions{}
	}

	if len(opt.AppVersion) == 0 {
		opt.AppVersion = os.Getenv("OUTPUT_APP_VERSION")
	}

	if len(opt.Levels) == 0 {
		opt.Levels = []logrus.Level{
			logrus.DebugLevel,
			logrus.TraceLevel,
		}
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
		stack: stackcache.New(6, "github.com/hatchify/output"),
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
	caller := h.stack.GetCaller()

	if len(caller.Function) > 0 {
		parts := strings.Split(caller.Function, "/")
		nameParts := strings.Split(parts[len(parts)-1], ".")
		e.Data["fn"] = nameParts[len(nameParts)-1]
	}

	callerFile := limitPath(caller.File, h.opt.PathSegmentsLimit)
	e.Data["src"] = fmt.Sprintf("%s:%d", callerFile, caller.Line)

	if len(h.opt.AppVersion) > 0 {
		e.Data["ver"] = h.opt.AppVersion
	}

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
