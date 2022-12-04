package sha

import (
	"path"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func Glob(s Sha, pc ...string) (globbed Sha, err error) {
	p := s.Path(pc...)

	var matches []string

	//TODO-P3 move to open_file_guard
	if matches, err = filepath.Glob(p + "*"); err != nil {
		err = errors.Wrap(err)
		return
	}

	if len(matches) == 0 {
		err = errors.Errorf("sha provided matches no objects: %s", p)
		return
	}

	if len(matches) > 1 {
		err = errors.Errorf(
			"ambiguous sha provided matches multiple objects: %q",
			matches,
		)

		return
	}

	p = string(matches[0])
	head := path.Base(path.Dir(p))
	tail := path.Base(p)

	if err = globbed.Set(head + tail); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
