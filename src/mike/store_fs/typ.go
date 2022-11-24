package store_fs

import (
	"fmt"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/golf/typ"
	"github.com/friedenberg/zit/src/hotel/cwd_files"
)

//TODO move to generics for store methods and combine types for all objekten
// type storeTyp struct {
// }

func (s *Store) CheckinTyp(p string) (t *typ.Named, err error) {
	return
}

func (s *Store) WriteTyp(t *typ.Named) (tco *typ.Typ, err error) {
	tcwd := &cwd_files.CwdTyp{
		FD: cwd_files.File{
			Path: fmt.Sprintf("%s.%s", t.Kennung, s.Konfig.Compiled.TypFileExtension),
		},
		Named: typ.Named{
			Kennung: t.Kennung,
		},
	}

	tco = &typ.Typ{
		External: *tcwd,
		Named:    *t,
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

func (s *Store) ReadTyp(ct *cwd_files.CwdTyp) (t *typ.Typ, err error) {
	format := typ.MakeFormatText(s.storeObjekten)

	var f *os.File

	if f, err = files.Open(ct.FD.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	t = &typ.Typ{
		External: *ct,
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
