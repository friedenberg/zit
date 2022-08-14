package open_file_guard

import (
	"os"
	"path/filepath"

	"github.com/friedenberg/zit/bravo/errors"
)

func DeleteFilesAndDirs(args ...string) (err error) {
	dirs := make(map[string]bool)

	for _, f := range args {
		if err = os.Remove(f); err != nil {
			err = errors.Error(err)
			return
		}

		// It's possible that the paths come in absolute or relative form. So we
		// convert each path into absolute to deduplicate and prevent trying to
		// remove the same directory more than once. That said, filepath.Abs does
		// not guarantee uniqueness, so it's still possible to experience an error.

		var abs string

		if abs, err = filepath.Abs(f); err != nil {
			err = errors.Error(err)
			return
		}

		dirs[filepath.Dir(abs)] = true
	}

	for d, _ := range dirs {
		var contents []string

		if contents, err = ReadDirNames(d); err != nil {
			err = errors.Error(err)
			return
		}

		if len(contents) == 0 {
			if err = os.Remove(d); err != nil {
				err = errors.Error(err)
				return
			}
		}
	}

	return
}
