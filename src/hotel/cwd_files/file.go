package cwd_files

import (
	"path"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/fd"
)

// TODO-P1 remove
type File = fd.FD

func MakeFile(dir string, p string) (ut File, err error) {
	ut = File{}

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
