package fd

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
	"github.com/friedenberg/zit/src/delta/thyme"
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
	fd.sha = sha.Make(akteWriter.GetShaLike())

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
		modTime: thyme.Tyme(fi.ModTime()),
	}

	if fd.path, err = filepath.Abs(path.Join(dir, fi.Name())); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeFileFromFD(
	fd FD,
	awf schnittstellen.AkteWriterFactory,
) (ut FD, err error) {
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

	ut.sha = sha.Make(aw.GetShaLike())

	return
}

func MakeFile(
	dir string,
	p string,
	awf schnittstellen.AkteWriterFactory,
) (ut FD, err error) {
	todo.Remove()
	ut = FD{}

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