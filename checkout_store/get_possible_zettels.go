package checkout_store

import (
	"os"
	"path"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/open_file_guard"
)

func (s Store) GetPossibleZettels() (hinweisen []string, err error) {
	hinweisen = make([]string, 0)

	var dirs []string

	if dirs, err = open_file_guard.ReadDirNames(s.path); err != nil {
		err = errors.Error(err)
		return
	}

	for _, d := range dirs {
		var fi os.FileInfo

		if fi, err = os.Stat(d); err != nil {
			err = errors.Error(err)
			return
		}

		if !fi.Mode().IsDir() {
			continue
		}

		var dirs2 []string

		if dirs2, err = open_file_guard.ReadDirNames(path.Join(s.path, d)); err != nil {
			err = errors.Error(err)
			return
		}

		for _, a := range dirs2 {
			if fi, err = os.Stat(path.Join(s.path, d, a)); err != nil {
				err = errors.Error(err)
				return
			}

			if fi.Mode().IsDir() {
				continue
			}

			if path.Ext(a) != ".md" {
				continue
			}

			//TODO hinweis validation?
			hinweisen = append(hinweisen, path.Join(d, a))
		}
	}

	return
}
