package umwelt

import (
	"io"
	"os"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/bravo/files"
	"github.com/friedenberg/zit/charlie/age"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/konfig"
)

type Umwelt struct {
	BasePath string
	cwd      string
	Konfig   konfig.Konfig
	Logger   stdprinter.Logger
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

	if u.cwd, err = os.Getwd(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (u Umwelt) Age() (a age.Age, err error) {
	fa := u.FileAge()

	if files.Exists(fa) {
		if a, err = age.Make(fa); err != nil {
			return
		}
	} else {
		a = age.MakeEmpty()
	}

	return
}

func (u Umwelt) DefaultEtiketten() (etiketten etikett.Set, err error) {
	etiketten = etikett.MakeSet()

	for e, t := range u.Konfig.Tags {
		if !t.AddToNewZettels {
			continue
		}

		if err = etiketten.AddString(e); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}
