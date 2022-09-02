package akten

import (
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type Akten interface {
	All() ([]sha.Sha, error)
}

type akten struct {
	basePath string
}

func New(basePath string) (s *akten, err error) {
	s = &akten{
		basePath: path.Join(basePath, "Objekte", "Akte"),
	}

	return
}

func (an *akten) All() (akte []sha.Sha, err error) {
	akte = make([]sha.Sha, 0)

	var dirs []string

	if dirs, err = open_file_guard.ReadDirNames(an.basePath); err != nil {
		err = errors.Error(err)
		return
	}

	for _, d := range dirs {
		var dirs2 []string

		if dirs2, err = open_file_guard.ReadDirNames(path.Join(an.basePath, d)); err != nil {
			err = errors.Error(err)
			return
		}

		for _, a := range dirs2 {
			var s sha.Sha

			if err = s.SetParts(d, a); err != nil {
				err = errors.Error(err)
				return
			}

			akte = append(akte, s)
		}
	}

	return
}
