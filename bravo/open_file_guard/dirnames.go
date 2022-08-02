package open_file_guard

import "github.com/friedenberg/zit/alfa/errors"

func ReadDirNames(p string) (names []string, err error) {
	d, err := Open(p)

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
