package blob

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
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
	BlobEnabledEnv    map[string]bool
}

// DefaultRetentionTTL is currently set to be 3 months.
const DefaultRenentionTTL = 2232 * time.Hour

func checkHookOptions(opt *HookOptions) *HookOptions {
	if opt == nil {
		opt = &HookOptions{}
	}

	if len(opt.Env) == 0 {
		opt.Env = os.Getenv("OUTPUT_ENV")
		if len(opt.Env) == 0 {
			opt.Env = "local"
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
		opt.BlobRetentionTTL = DefaultRenentionTTL
	}

	if len(opt.BlobEnabledEnv) == 0 {
		opt.BlobEnabledEnv = map[string]bool{
			"prod":    true,
			"staging": true,
			"test":    true,
		}
	}

	return opt
}

// NewHook initializes a new output.Hook using provided params and options.
func NewHook(opt *HookOptions) (blobHook logrus.Hook, err error) {
	h := new(hook)
	h.opt = checkHookOptions(opt)

	out := logrus.WithFields(logrus.Fields{
		"account":  opt.BlobStoreAccount,
		"bucket":   opt.BlobStoreBucket,
		"endpoint": opt.BlobStoreEndpoint,
	})

	if h.s3Remote, err = NewS3Remote(
		opt.BlobStoreAccount,
		opt.BlobStoreKey,
		opt.BlobStoreEndpoint,
		opt.BlobStoreRegion,
		opt.BlobStoreBucket,
	); err != nil {
		out.WithError(err).Warning("failed to init S3 session")

		h.s3Remote = nil

		return h, err
	} else if err := h.s3Remote.CheckAccess(opt.Env); err != nil {
		out.WithError(err).Warning("failed to verify S3 remote access")

		h.s3Remote = nil

		return h, err
	}

	return h, nil
}

type hook struct {
	opt      *HookOptions
	s3Remote S3Remote
}

func (h *hook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.TraceLevel,
	}
}

func (h *hook) Fire(e *logrus.Entry) error {
	blob, hasBlob := e.Data["blob"]
	if !hasBlob {
		return nil
	}

	if h.s3Remote == nil {
		logrus.Warning("blob provided but S3 remote is disabled")
		delete(e.Data, "blob")

		return nil
	} else if enabled := h.opt.BlobEnabledEnv[h.opt.Env]; !enabled {
		logrus.Infof("blob provided but uploading is disabled in %s", h.opt.Env)
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
		logrus.WithError(err).WithFields(logrus.Fields{
			"bucket": h.opt.BlobStoreBucket,
			"key":    objectKey,
		}).Errorln("failed to upload blob to S3 remote server")
	}
}
