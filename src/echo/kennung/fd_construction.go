package kennung

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/sha"
)

func FDFromDir(
	p string,
) (fd FD, err error) {
	fd.isDir = true
	fd.path = p

	return
}

func FDFromPathWithAkteWriterFactory(
	p string,
	awf schnittstellen.AkteWriterFactory,
) (fd FD, err error) {
	if p == "" {
		err = errors.Errorf("empty path")
		return
	}

	if awf == nil {
		panic("schnittstellen.AkteWriterFactory is nil")
	}

	var f *os.File

	if f, err = files.OpenExclusiveReadOnly(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	var akteWriter sha.WriteCloser

	if akteWriter, err = awf.AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, akteWriter)

	if _, err = io.Copy(akteWriter, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	var fi os.FileInfo

	if fi, err = f.Stat(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fd, err = FileInfo(fi, path.Dir(p)); err != nil {
		err = errors.Wrap(err)
		return
	}

	fd.path = p
	fd.Sha = sha.Make(akteWriter.GetShaLike())

	return
}

func FDFromPath(p string) (fd FD, err error) {
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

	if fd, err = FileInfo(fi, path.Dir(p)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func File(f *os.File) (fd FD, err error) {
	if f == nil {
		err = errors.Errorf("nil file desriptor")
		return
	}

	var fi os.FileInfo

	if fi, err = f.Stat(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if fd, err = FileInfo(fi, path.Dir(f.Name())); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func FileInfo(fi os.FileInfo, dir string) (fd FD, err error) {
	fd = FD{
		isDir:   fi.IsDir(),
		ModTime: Tyme(fi.ModTime()),
	}

	if fd.path, err = filepath.Abs(path.Join(dir, fi.Name())); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
