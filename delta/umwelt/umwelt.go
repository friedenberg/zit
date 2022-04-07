package umwelt

import (
	"fmt"
	"io"
	"os"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/charlie/konfig"
)

type Umwelt struct {
	BasePath string
	cwd      string
	Konfig   konfig.Konfig
	Logger   _Logger
	In       io.Reader
	Out      io.Writer
	Err      io.Writer
}

func MakeUmwelt(c konfig.Konfig) (u *Umwelt, err error) {
	u = &Umwelt{
		Konfig: c,
		Logger: c.Logger,
		In:     os.Stdin,
		Out:    os.Stdout,
		Err:    os.Stderr,
	}

	if u.BasePath, err = c.DirZit(); err != nil {
		err = errors.Error(err)
		return
	}

	fmt.Println(u.DirZit())

	if u.cwd, err = os.Getwd(); err != nil {
		err = errors.Error(err)
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
