package cwd

import (
	"io"
	"os"
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/fd"
)

func MakeFileFromFD(
	fd fd.FD,
	awf schnittstellen.AkteWriterFactory,
) (ut fd.FD, err error) {
	ut = fd

	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(ut.GetPath()); err != nil {
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
) (ut fd.FD, err error) {
	todo.Remove()
	ut = fd.FD{}

	p = path.Join(dir, p)

	if err = ut.Set(p); err != nil {
		err = errors.Wrapf(err, "path: %q", p)
		return
	}

	if err = ut.SetPathRel(ut.GetPath(), dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ut, err = MakeFileFromFD(ut, awf); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
