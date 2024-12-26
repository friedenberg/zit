package dir_layout

import (
	"context"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

// TODO only call reset temp when actually not resetting temp
func (s Layout) ResetTempOnExit(ctx errors.Context) (err error) {
	errIn := context.Cause(ctx)

	if errIn != nil || s.GetDebug().NoTempDirCleanup {
		// ui.Err().Printf("temp dir: %q", s.DirTempLocal())
	} else {
		if err = os.RemoveAll(s.TempLocal.BasePath); err != nil {
			err = errors.Wrapf(err, "failed to remove temp local")
			return
		}
	}

	return
}

type TemporaryFS struct {
	BasePath string
}

func (s TemporaryFS) DirTemp() (d string, err error) {
	if d, err = os.MkdirTemp(s.BasePath, ""); err != nil {
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
	if f, err = os.CreateTemp(s.BasePath, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
