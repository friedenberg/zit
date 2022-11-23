package store_fs

import (
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/echo/typ"
	"github.com/friedenberg/zit/src/foxtrot/cwd_files"
	"github.com/friedenberg/zit/src/foxtrot/typ_checked_out"
)

func (s *Store) CheckinTyp(p string) (t *typ.Typ, err error) {
	return
}

func (s *Store) ReadTyp(ct *cwd_files.CwdTyp) (t *typ_checked_out.Typ, err error) {
	format := typ.MakeFormatText(s.storeObjekten)

	var f *os.File

	if f, err = files.Open(ct.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	t = &typ_checked_out.Typ{
		Path: ct.Path,
		Typ: typ.Typ{
			Kennung: ct.Kennung,
		},
	}

	if _, err = format.ReadFormat(f, &t.Typ); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
