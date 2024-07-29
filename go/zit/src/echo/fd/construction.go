package fd

import (
	"io"
	"os"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

func MakeFromDirPath(
	p string,
) (fd *FD, err error) {
	fd = &FD{}
	fd.isDir = true
	fd.path = p

	return
}

func MakeFromPath(p string) (fd *FD, err error) {
	if p == "" {
		err = errors.Errorf("nil file desriptor")
		return
	}

	if p == "." {
		err = errors.Errorf("'.' not supported")
		return
	}

	var fi os.FileInfo

	if fi, err = os.Stat(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fd, err = MakeFromFileInfoWithDir(fi, path.Dir(p)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeFromFileInfoWithDir(fi os.FileInfo, dir string) (fd *FD, err error) {
	fd = &FD{}
	err = fd.SetFileInfoWithDir(fi, dir)
	return
}

func MakeFromFileFromFD(
	fd *FD,
	awf interfaces.BlobWriterFactory,
) (ut *FD, err error) {
	ut = &FD{}
	ut.ResetWith(fd)

	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(ut.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	var aw sha.WriteCloser

	if aw, err = awf.BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if _, err = io.Copy(aw, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	ut.sha.SetShaLike(aw)

	return
}

func MakeFromPathWithBlobWriterFactory(
	p string,
	awf interfaces.BlobWriterFactory,
) (ut *FD, err error) {
	todo.Remove()
	ut = &FD{}

	if err = ut.Set(p); err != nil {
		err = errors.Wrapf(err, "path: %q", p)
		return
	}

	if ut, err = MakeFromFileFromFD(ut, awf); err != nil {
		err = errors.Wrapf(err, "Path: %q", p)
		return
	}

	return
}
