package umwelt

import (
	"io"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/delta/standort"
)

type Umwelt struct {
	BasePath string
	cwd      string
	standort.Standort
	Konfig konfig.Konfig
	Logger errors.Logger
	In     io.Reader
	Out    io.Writer
	Err    io.Writer
}

func MakeUmwelt(c konfig.Konfig) (u *Umwelt, err error) {
	u = &Umwelt{
		Konfig: c,
		Logger: c.Logger,
		In:     os.Stdin,
		Out:    os.Stdout,
		Err:    os.Stderr,
	}

	if u.Standort, err = standort.Make(c); err != nil {
		err = errors.Wrap(err)
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
			err = errors.Wrap(err)
			return
		}
	}

	return
}
