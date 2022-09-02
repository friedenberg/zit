package store_working_directory

import (
	"os"
	"path"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
)

type CwdFiles struct {
	Zettelen         []string
	Akten            []string
	EmptyDirectories []string
}

func (c CwdFiles) Len() int {
	return len(c.Zettelen) + len(c.Akten)
}

func (s Store) GetPossibleZettels() (result CwdFiles, err error) {
	result.Zettelen = make([]string, 0)
	result.Akten = make([]string, 0)

	var dirs []string

	if dirs, err = open_file_guard.ReadDirNames(s.path); err != nil {
		err = errors.Error(err)
		return
	}

	for _, d := range dirs {
		if strings.HasPrefix(d, ".") {
			continue
		}

		d2 := path.Join(s.path, d)

		var fi os.FileInfo

		if fi, err = os.Stat(d); err != nil {
			err = errors.Error(err)
			return
		}

		if !fi.Mode().IsDir() {
			continue
		}

		var dirs2 []string

		if dirs2, err = open_file_guard.ReadDirNames(d2); err != nil {
			err = errors.Error(err)
			return
		}

		if len(dirs2) == 0 {
			result.EmptyDirectories = append(result.EmptyDirectories, d2)
		}

		for _, a := range dirs2 {
			if strings.HasPrefix(a, ".") {
				continue
			}

			if fi, err = os.Stat(path.Join(s.path, d, a)); err != nil {
				err = errors.Error(err)
				return
			}

			if fi.Mode().IsDir() {
				continue
			}

			p := path.Join(d, a)
			//TODO-decision: should there be hinweis validation?

			//TODO-refactor: akten vs zettel file extensions
			if path.Ext(a) == ".md" {
				result.Zettelen = append(result.Zettelen, p)
			} else {
				result.Akten = append(result.Akten, p)
			}
		}
	}

	return
}
