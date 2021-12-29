package sha

import (
	"path"
	"path/filepath"
)

func (s Sha) Glob(pc ...string) (globbed Sha, err error) {
	p := s.Path(pc...)

	var matches []string

	//TODO move to open_file_guard
	if matches, err = filepath.Glob(p + "*"); err != nil {
		err = _Error(err)
		return
	}

	if len(matches) == 0 {
		err = _Errorf("sha provided matches no objects: %s", p)
		return
	}

	if len(matches) > 1 {
		err = _Errorf(
			"ambiguous sha provided matches multiple objects: %q",
			matches,
		)

		return
	}

	p = string(matches[0])
	head := path.Base(path.Dir(p))
	tail := path.Base(p)

	if err = globbed.Set(head + tail); err != nil {
		err = _Error(err)
		return
	}

	return
}
