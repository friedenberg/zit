package open_file_guard

import (
	"os"
	"path"
)

func DeleteFilesAndDirs(args ...string) (err error) {
	dirs := make(map[string]bool)

	for _, f := range args {
		if err = os.Remove(f); err != nil {
			err = _Error(err)
			return
		}

		dirs[path.Dir(f)] = true
	}

	for d, _ := range dirs {
		var contents []string

		if contents, err = ReadDirNames(d); err != nil {
			err = _Error(err)
			return
		}

		if len(contents) == 0 {
			if err = os.Remove(d); err != nil {
				err = _Error(err)
				return
			}
		}
	}

	return
}
