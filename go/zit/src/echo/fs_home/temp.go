package fs_home

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

// TODO only call reset temp when actually not resetting temp
func (s Home) ResetTempOnExit(errIn error) (err error) {
	if errIn != nil || s.debug.NoTempDirCleanup {
		// ui.Err().Printf("temp dir: %q", s.DirTempLocal())
	} else {
		if err = os.RemoveAll(s.TempLocal.basePath); err != nil {
			err = errors.Wrapf(err, "failed to remove temp local")
			return
		}
	}

	return
}

type TemporaryFS struct {
	basePath string
}

func (s TemporaryFS) DirTemp() (d string, err error) {
	if d, err = os.MkdirTemp(s.basePath, ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s TemporaryFS) FileTemp() (f *os.File, err error) {
	if f, err = s.FileTempWithTemplate(""); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s TemporaryFS) FileTempWithTemplate(t string) (f *os.File, err error) {
	if f, err = os.CreateTemp(s.basePath, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
