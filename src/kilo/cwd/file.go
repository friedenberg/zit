package cwd

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
)

func MakeFileFromFD(
	fd kennung.FD,
	awf schnittstellen.AkteWriterFactory,
) (ut kennung.FD, err error) {
	ut = fd

	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(ut.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	var aw sha.WriteCloser

	if aw, err = awf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if _, err = io.Copy(aw, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	ut.Sha = sha.Make(aw.GetShaLike())

	return
}

func MakeFile(
	dir string,
	p string,
	awf schnittstellen.AkteWriterFactory,
) (ut kennung.FD, err error) {
	todo.Remove()
	ut = kennung.FD{}

	p = path.Join(dir, p)

	if err = ut.Set(p); err != nil {
		err = errors.Wrapf(err, "path: %q", p)
		return
	}

	if ut.Path, err = filepath.Rel(dir, ut.Path); err != nil {
		err = errors.Wrapf(err, "path: %q", ut.Path)
		return
	}

	if ut, err = MakeFileFromFD(ut, awf); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
