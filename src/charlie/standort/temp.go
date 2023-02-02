package standort

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func (s Standort) ResetTemp() (err error) {
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
	return s.DirZit("tmp")
}

func (s Standort) FileTempOS() (f *os.File, err error) {
	if f, err = os.CreateTemp("", ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s Standort) FileTempLocal() (f *os.File, err error) {
	if f, err = os.CreateTemp(s.DirTempLocal(), ""); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
