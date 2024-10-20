package fs_home

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"syscall"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

func (s Home) DeleteAll(p string) (err error) {
	if s.dryRun {
		return
	}

	if err = os.RemoveAll(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Home) Delete(p string) (err error) {
	if s.dryRun {
		return
	}

	if err = os.Remove(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO only call reset temp when actually not resetting temp
func (s Home) ResetTempOnExit(errIn error) (err error) {
	if errIn != nil || s.debug.NoTempDirCleanup {
		// ui.Err().Printf("temp dir: %q", s.DirTempLocal())
	} else {
		if err = os.RemoveAll(s.DirTempLocal()); err != nil {
			err = errors.Wrapf(err, "failed to remove temp local")
			return
		}
	}

	return
}

func (s Home) DirTempOS() (d string, err error) {
	if d, err = os.MkdirTemp("", ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Home) DirTempLocal() string {
	return s.DirZit(fmt.Sprintf("tmp-%d", s.pid))
}

func (s Home) FileTempOS() (f *os.File, err error) {
	if f, err = os.CreateTemp("", ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Home) FileTempLocal() (f *os.File, err error) {
	if f, err = s.FileTempLocalWithTemplate(""); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Home) FileTempLocalWithTemplate(t string) (f *os.File, err error) {
	if f, err = os.CreateTemp(s.DirTempLocal(), t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Home) FifoPipeWithExtension(ext string) (p string, err error) {
	p = path.Join(
		s.DirTempLocal(),
		fmt.Sprintf("%s.%s", strconv.Itoa(rand.Int()), ext),
	)

	if err = syscall.Mknod(p, syscall.S_IFIFO|0o666, 0); err != nil {
		err = errors.Wrapf(err, "Path: %s", p)
		return
	}

	// if err = os.Remove(p); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	return
}

func (s Home) FifoPipe() (p string, err error) {
	p = path.Join(
		s.DirTempLocal(),
		strconv.Itoa(rand.Int()),
	)

	if err = syscall.Mknod(p, syscall.S_IFIFO|0o666, 0); err != nil {
		err = errors.Wrapf(err, "Path: %s", p)
		return
	}

	// if err = os.Remove(p); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	return
}
