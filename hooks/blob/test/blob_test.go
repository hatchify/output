package blob

import (
	"os"
	"testing"
	"time"

	blobHook "github.com/hatchify/output/hooks/blob"

	"github.com/hatchify/output"
)

func TestBlobHook(t *testing.T) {
	testBlob := []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna
aliqua. Ut enim ad minim veniam, quis nostrud exercitation
ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit
esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
occaecat cupidatat non proident, sunt in culpa qui officia
deserunt mollit anim id est laborum.`)

	opts := &blobHook.HookOptions{
		Env: "test",
	}

	hook, _ := blobHook.NewHook(opts)
	out := output.NewOutputter(os.Stderr, new(output.TextFormatter), hook)
	ts := time.Now()

	out.WithField("blob", testBlob).Infoln("test is running, trying to submit blob")
	out.Debug("test done in %s", time.Since(ts))
}
