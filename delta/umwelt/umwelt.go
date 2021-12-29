package umwelt

import (
	"io"
	"os"
	"path"
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

func (u Umwelt) Dir() string {
	return u.BasePath
}

func (u Umwelt) DirZit(p ...string) string {
	return path.Join(
		append(
			[]string{u.Dir(), ".zit"},
			p...,
		)...,
	)
}

func (u Umwelt) Age() (a _Age, err error) {
	p := u.DirZit()

	if a, err = _AgeMake(p); err != nil {
		return
	}

	return
}
