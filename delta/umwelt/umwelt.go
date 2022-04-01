package umwelt

import (
	"io"
	"os"
)

type Umwelt struct {
	BasePath string
	cwd      string
	Konfig   _Konfig
	Logger   _Logger
	In       io.Reader
	Out      io.Writer
	Err      io.Writer
}

func MakeUmwelt(c _Konfig) (u *Umwelt, err error) {
	u = &Umwelt{
		Konfig: c,
		Logger: c.Logger,
		In:     os.Stdin,
		Out:    os.Stdout,
		Err:    os.Stderr,
	}

	if u.BasePath, err = c.DirZit(); err != nil {
		err = _Error(err)
		return
	}

	if u.cwd, err = os.Getwd(); err != nil {
		err = _Error(err)
		return
	}

	return
}

func (u Umwelt) Age() (a _Age, err error) {
	fa := u.FileAge()

	if _FilesExist(fa) {
		if a, err = _AgeMake(fa); err != nil {
			return
		}
	} else {
		a = _AgeMakeEmpty()
	}

	return
}
