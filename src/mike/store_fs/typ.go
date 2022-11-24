package store_fs

import (
	"fmt"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/echo/typ"
	"github.com/friedenberg/zit/src/foxtrot/cwd_files"
	"github.com/friedenberg/zit/src/foxtrot/typ_checked_out"
)

//TODO move to generics for store methods and combine types for all objekten
// type storeTyp struct {
// }

func (s *Store) CheckinTyp(p string) (t *typ.Named, err error) {
	return
}

func (s *Store) WriteTyp(t *typ.Named) (tco *typ_checked_out.Typ, err error) {
	tcwd := &cwd_files.CwdTyp{
		FD: cwd_files.File{
			Path: fmt.Sprintf("%s.%s", t.Kennung, s.Konfig.Compiled.TypFileExtension),
		},
		Named: typ.Named{
			Kennung: t.Kennung,
		},
	}

	tco = &typ_checked_out.Typ{
		CwdTyp: *tcwd,
		Named:  *t,
	}

	var f *os.File

	if f, err = files.CreateExclusiveWriteOnly(tcwd.FD.Path); err != nil {
		if errors.IsExist(err) {
			tco, err = s.ReadTyp(tcwd)
		} else {
			err = errors.Wrap(err)
		}
		return
	}

	defer errors.Deferred(&err, f.Close)

	format := typ.MakeFormatText(s.storeObjekten)

	if _, err = format.WriteFormat(f, &t.Stored.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadTyp(ct *cwd_files.CwdTyp) (t *typ_checked_out.Typ, err error) {
	format := typ.MakeFormatText(s.storeObjekten)

	var f *os.File

	if f, err = files.Open(ct.FD.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	t = &typ_checked_out.Typ{
		CwdTyp: *ct,
		Named: typ.Named{
			Kennung: ct.Named.Kennung,
		},
	}

	if _, err = format.ReadFormat(f, &t.Named.Stored.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
