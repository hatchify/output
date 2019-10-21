package blob

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hatchify/output"
)

// HookOptions allows to set additional Hook options.
type HookOptions struct {
	Env               string
	BlobStoreURL      string
	BlobStoreAccount  string
	BlobStoreKey      string
	BlobStoreEndpoint string
	BlobStoreRegion   string
	BlobStoreBucket   string
	BlobRetentionTTL  time.Duration
}

func checkHookOptions(opt *HookOptions) *HookOptions {
	if opt == nil {
		opt = &HookOptions{}
	}
	if len(opt.Env) == 0 {
		opt.Env = os.Getenv("OUTPUT_ENV")
		if len(opt.Env) == 0 {
			opt.Env = "test"
		}
	}
	if len(opt.BlobStoreURL) == 0 {
		opt.BlobStoreURL = os.Getenv("OUTPUT_BLOB_STORE_URL")
	}
	if len(opt.BlobStoreAccount) == 0 {
		opt.BlobStoreAccount = os.Getenv("OUTPUT_BLOB_STORE_ACCOUNT")
	}
	if len(opt.BlobStoreKey) == 0 {
		opt.BlobStoreKey = os.Getenv("OUTPUT_BLOB_STORE_KEY")
	}
	if len(opt.BlobStoreEndpoint) == 0 {
		opt.BlobStoreEndpoint = os.Getenv("OUTPUT_BLOB_STORE_ENDPOINT")
	}
	if len(opt.BlobStoreRegion) == 0 {
		opt.BlobStoreRegion = os.Getenv("OUTPUT_BLOB_STORE_REGION")
	}
	if len(opt.BlobStoreBucket) == 0 {
		opt.BlobStoreBucket = os.Getenv("OUTPUT_BLOB_STORE_BUCKET")
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
	s3Remote := NewS3Remote(
		opt.BlobStoreAccount,
		opt.BlobStoreKey,
		opt.BlobStoreEndpoint,
		opt.BlobStoreRegion,
		opt.BlobStoreBucket,
	)
	if err := s3Remote.CheckAccess(opt.Env); err != nil {
		output.WithError(err).WithFields(output.Fields{
			"account": opt.BlobStoreAccount,
			"bucket":  opt.BlobStoreBucket,
			"endpoint": opt.BlobStoreEndpoint,
		}).Warning("failed to verify S3 remote access")
		s3Remote = nil
	}
	return &hook{
		opt:      opt,
		s3Remote: s3Remote,
	}
}

type hook struct {
	opt      *HookOptions
	s3Remote S3Remote
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
	blob, hasBlob := e.Data["blob"]
	if !hasBlob {
		return nil
	}
	if h.s3Remote == nil {
		output.Warning("blob provided but S3 remote is disabled")
		delete(e.Data, "blob")
		return nil
	}
	var blobPayload []byte
	switch bb := blob.(type) {
	case string:
		blobPayload = []byte(bb)
	case []byte:
		blobPayload = make([]byte, len(bb))
		copy(blobPayload, bb)
	default:
		delete(e.Data, "blob")
		return nil
	}
	blobID := NewBlobID()
	if len(h.opt.BlobStoreURL) > 0 {
		e.Data["blob"] = fmt.Sprintf("%s/%s", h.opt.BlobStoreURL, blobID)
	} else {
		e.Data["blob"] = fmt.Sprintf("%s/%s", h.opt.Env, blobID)
	}
	h.blobUpload(blobID, blobPayload)
	return nil
}

func (h *hook) blobUpload(blobID string, payload []byte) {
	objectKey := filepath.Join(h.opt.Env, blobID)
	_, err := h.s3Remote.PutObject(objectKey, bytes.NewReader(payload), nil)
	if err != nil {
		output.WithError(err).WithFields(output.Fields{
			"bucket": h.opt.BlobStoreBucket,
			"key":    objectKey,
		}).Errorln("failed to upload blob to S3 remote server")
	}
}