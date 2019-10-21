package s3upload

import (
	"time"

	"github.com/hatchify/output"
)

// HookOptions allows to set additional Hook options.
type HookOptions struct {
	BlobRetentionTTL time.Duration
}

func checkHookOptions(opt *HookOptions) *HookOptions {
	if opt == nil {
		opt = &HookOptions{}
	}
	if opt.BlobRetentionTTL == 0 {
		// keep blobs for 3 months
		opt.BlobRetentionTTL = 2232 * time.Hour
	}
	return opt
}

// NewHook initializes a new output.Hook using provided params and options.
func NewHook(opt *HookOptions) output.Hook {
	opt = checkHookOptions(opt)
	return &hook{
		opt: opt,
	}
}

type hook struct {
	opt *HookOptions
}

func (h *hook) Levels() []output.Level {
	return []output.Level{
		output.PanicLevel,
		output.FatalLevel,
		output.ErrorLevel,
		output.WarnLevel,
		output.SuccessLevel,
		output.InfoLevel,
		output.DebugLevel,
		output.TraceLevel,
	}
}

func (h *hook) Fire(e *output.Entry) error {
	return nil
}
