package umwelt

import (
	"io"
	"os"
)

type Umwelt struct {
	BasePath string
	Lock     *_FileLock
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

	u.Lock = _FileLockNew(u.DirZit("Lock"))

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
