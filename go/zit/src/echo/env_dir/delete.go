package env_dir

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

func (env env) Delete(paths ...string) (err error) {
	for _, path := range paths {
		path = filepath.Clean(path)

		if path == "." {
			err = errors.Errorf("invalid delete request: %q", path)
			return
		}

		if env.IsDryRun() {
			ui.Err().Print("would delete:", path)
			return
		}

		if err = os.RemoveAll(path); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
