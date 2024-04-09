package standort

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"syscall"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

func (s Standort) ResetTemp() (err error) {
	if s.debug.NoTempDirCleanup {
		return
	}

	if err = os.RemoveAll(s.DirTempLocal()); err != nil {
		err = errors.Wrapf(err, "failed to remove temp local")
		return
	}

	return
}

func (s Standort) DirTempOS() (d string, err error) {
	if d, err = os.MkdirTemp("", ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Standort) DirTempLocal() string {
	return s.DirZit(fmt.Sprintf("tmp-%d", s.pid))
}

func (s Standort) FileTempOS() (f *os.File, err error) {
	if f, err = os.CreateTemp("", ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Standort) FileTempLocal() (f *os.File, err error) {
	if f, err = s.FileTempLocalWithTemplate(""); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Standort) FileTempLocalWithTemplate(t string) (f *os.File, err error) {
	if f, err = os.CreateTemp(s.DirTempLocal(), t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Standort) FifoPipe() (p string, err error) {
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
