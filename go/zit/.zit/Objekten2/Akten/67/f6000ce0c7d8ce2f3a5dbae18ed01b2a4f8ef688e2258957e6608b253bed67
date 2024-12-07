package fd

import (
	"io"
	"io/fs"
	"os"
	"path"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
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

func MakeFromPathAndDirEntry(
	p string,
	de fs.DirEntry,
	awf interfaces.BlobWriterFactory,
) (fd *FD, err error) {
	if p == "" {
		err = errors.Errorf("nil file desriptor")
		return
	}

	if p == "." {
		err = errors.Errorf("'.' not supported")
		return
	}

	var fi os.FileInfo

	if fi, err = de.Info(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fd, err = MakeFromFileInfoWithDir(fi, path.Dir(p), awf); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeFromPath(
	p string,
	awf interfaces.BlobWriterFactory,
) (fd *FD, err error) {
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

	if fd, err = MakeFromFileInfoWithDir(
		fi,
		path.Dir(p),
		awf,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeFromFileInfoWithDir(
	fi os.FileInfo,
	dir string,
	awf interfaces.BlobWriterFactory,
) (fd *FD, err error) {
	// TODO use pool
	fd = &FD{}

	if err = fd.SetFileInfoWithDir(fi, dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fi.IsDir() {
		return
	}

	// TODO eventually enforce requirement of blob writer factory
	if awf == nil {
		return
	}

	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(fd.GetPath()); err != nil {
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

	fd.sha.SetShaLike(aw)
	fd.state = StateStored

	return
}
