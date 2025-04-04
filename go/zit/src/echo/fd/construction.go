package fd

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

func MakeFromDirPath(
	path string,
) (fd *FD, err error) {
	fd = &FD{}
	fd.isDir = true
	fd.path = path

	return
}

func MakeFromPathAndDirEntry(
	path string,
	dirEntry fs.DirEntry,
	blobWriter interfaces.BlobWriter,
) (fd *FD, err error) {
	if path == "" {
		err = errors.Errorf("nil file desriptor")
		return
	}

	if path == "." {
		err = errors.Errorf("'.' not supported")
		return
	}

	var fi os.FileInfo

	if fi, err = dirEntry.Info(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fd, err = MakeFromFileInfoWithDir(fi, filepath.Dir(path), blobWriter); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeFromPath(
	baseDir string,
	path string,
	blobWriter interfaces.BlobWriter,
) (fd *FD, err error) {
	if path == "" {
		err = errors.Errorf("nil file desriptor")
		return
	}

	if path == "." {
		err = errors.Errorf("'.' not supported")
		return
	}

	if !filepath.IsAbs(path) {
		path = filepath.Clean(filepath.Join(baseDir, path))
	}

	var fi os.FileInfo

	if fi, err = os.Stat(path); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fd, err = MakeFromFileInfoWithDir(
		fi,
		filepath.Dir(path),
		blobWriter,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeFromFileInfoWithDir(
	fileInfo os.FileInfo,
	dir string,
	blobWriter interfaces.BlobWriter,
) (fd *FD, err error) {
	// TODO use pool
	fd = &FD{}

	if err = fd.SetFileInfoWithDir(fileInfo, dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fileInfo.IsDir() {
		return
	}

	// TODO eventually enforce requirement of blob writer factory
	if blobWriter == nil {
		return
	}

	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(fd.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	var aw sha.WriteCloser

	if aw, err = blobWriter.BlobWriter(); err != nil {
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
