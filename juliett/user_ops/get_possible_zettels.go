package user_ops

import (
	"os"
	"path"

	"github.com/friedenberg/zit/bravo/open_file_guard"
	"github.com/friedenberg/zit/delta/umwelt"
)

type GetPossibleZettels struct {
	Umwelt *umwelt.Umwelt
}

type GetPossibleZettelsResults struct {
	Hinweisen []string
}

func (c GetPossibleZettels) Run() (results GetPossibleZettelsResults, err error) {
	results.Hinweisen = make([]string, 0)

	var cwd string

	if cwd, err = os.Getwd(); err != nil {
		err = _Error(err)
		return
	}

	var dirs []string

	if dirs, err = open_file_guard.ReadDirNames(cwd); err != nil {
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

		if dirs2, err = open_file_guard.ReadDirNames(path.Join(cwd, d)); err != nil {
			err = _Error(err)
			return
		}

		for _, a := range dirs2 {
			if fi, err = os.Stat(path.Join(cwd, d, a)); err != nil {
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
			results.Hinweisen = append(results.Hinweisen, path.Join(d, a))
		}
	}

	return

	return
}
