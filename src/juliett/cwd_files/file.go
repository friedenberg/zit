package cwd_files

import (
	"path"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/fd"
)

func MakeFile(dir string, p string) (ut fd.FD, err error) {
	ut = fd.FD{}

	p = path.Join(dir, p)

	if err = ut.Set(p); err != nil {
		err = errors.Wrapf(err, "path: %q", p)
		return
	}

	if ut.Path, err = filepath.Rel(dir, ut.Path); err != nil {
		err = errors.Wrapf(err, "path: %q", ut.Path)
		return
	}

	return
}
