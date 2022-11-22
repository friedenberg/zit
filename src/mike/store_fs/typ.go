package store_fs

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/echo/typ"
)

func (s *Store) CheckinTyp(p string) (t *typ.Typ, err error) {
	return
}

func (s *Store) CheckoutTyp(p string) (t *typ.Typ, err error) {
	return
}

func (s *Store) ReadTyp(p string) (t *typ.Typ, err error) {
	format := typ.MakeFormatText(s.storeObjekten)

	var f *os.File

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	t = &typ.Typ{}

	if _, err = format.ReadFormat(f, t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
