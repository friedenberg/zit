package files

import (
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
)

func RemoveIfExists(p string) (err error) {
	err = os.Remove(p)

	if errors.IsNotExist(err) {
		err = nil
	}

	return
}

func DeleteFilesAndDirsSet(
	fs schnittstellen.SetLike[schnittstellen.Stringer],
) (err error) {
	return fs.Each(
		func(f schnittstellen.Stringer) (err error) {
			if err = os.Remove(f.String()); err != nil {
				if errors.IsNotExist(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	)
}

func DeleteFilesAndDirs(args ...string) (err error) {
	dirs := make(map[string]bool)

	for _, f := range args {
		if err = os.Remove(f); err != nil {
			if errors.IsNotExist(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		// It's possible that the paths come in absolute or relative form. So we
		// convert each path into absolute to deduplicate and prevent trying to
		// remove the same directory more than once. That said, filepath.Abs
		// does not guarantee uniqueness, so it's still possible to experience
		// an error.

		var abs string

		if abs, err = filepath.Abs(f); err != nil {
			err = errors.Wrap(err)
			return
		}

		dirs[filepath.Dir(abs)] = true
	}

	for d := range dirs {
		var contents []string

		if contents, err = ReadDirNames(d); err != nil {
			err = errors.Wrap(err)
			return
		}

		if len(contents) == 0 {
			if err = os.Remove(d); err != nil {
				if errors.IsNotExist(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
					return
				}
			}
		}
	}

	return
}
