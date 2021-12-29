package zettels

import (
	"os"
	"path"
)

func (zs *zettels) GetPossibleZettels(wd string) (hins []string, err error) {
	hins = make([]string, 0)

	var dirs []string

	if dirs, err = _ReadDirNames(wd); err != nil {
		err = _Error(err)
		return
	}

	for _, d := range dirs {
		var fi os.FileInfo

		if fi, err = os.Stat(d); err != nil {
			err = _Error(err)
			return
		}

		if !fi.Mode().IsDir() {
			continue
		}

		var dirs2 []string

		if dirs2, err = _ReadDirNames(path.Join(wd, d)); err != nil {
			err = _Error(err)
			return
		}

		for _, a := range dirs2 {
			if fi, err = os.Stat(path.Join(wd, d, a)); err != nil {
				err = _Error(err)
				return
			}

			if fi.Mode().IsDir() {
				continue
			}

			if path.Ext(a) != ".md" {
				continue
			}

			//TODO hinweis validation?
			hins = append(hins, path.Join(d, a))
		}
	}

	return
}
