package akten

import "path"

type Akten interface {
	All() ([]_Sha, error)
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

func (an *akten) All() (akte []_Sha, err error) {
	akte = make([]_Sha, 0)

	var dirs []string

	if dirs, err = _ReadDirNames(an.basePath); err != nil {
		err = _Error(err)
		return
	}

	for _, d := range dirs {
		var dirs2 []string

		if dirs2, err = _ReadDirNames(path.Join(an.basePath, d)); err != nil {
			err = _Error(err)
			return
		}

		for _, a := range dirs2 {
			var s _Sha

			if err = s.SetParts(d, a); err != nil {
				err = _Error(err)
				return
			}

			akte = append(akte, s)
		}
	}

	return
}
