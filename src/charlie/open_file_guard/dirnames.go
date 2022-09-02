package open_file_guard

import (
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func ReadDirNames(ps ...string) (names []string, err error) {
	d, err := Open(path.Join(ps...))

	if err != nil {
		err = errors.Error(err)
		return
	}

	defer Close(d)

	if names, err = d.Readdirnames(0); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
