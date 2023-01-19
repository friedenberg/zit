package store_fs

import (
	"fmt"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/fd"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/hotel/typ"
)

func (s *Store) CheckinTyp(p string) (t *typ.Transacted, err error) {
	return
}

func (s *Store) WriteTyp(t *typ.Transacted) (te *typ.External, err error) {
	te = &typ.External{
		FD: fd.FD{
			Path: fmt.Sprintf("%s.%s", t.Kennung(), s.erworben.FileExtensions.Typ),
		},
		//TODO-P2 move to central place
		Objekte: t.Objekte,
		Sku: sku.External[kennung.Typ, *kennung.Typ]{
			ObjekteSha: sha.Make(t.GetObjekteSha()),
			Kennung:    t.Sku.Kennung,
		},
	}

	var f *os.File

	if f, err = files.CreateExclusiveWriteOnly(te.FD.Path); err != nil {
		if errors.IsExist(err) {
			err = s.ReadTyp(te)
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.Deferred(&err, f.Close)

	format := typ.MakeFormatText(s.storeObjekten)

	if _, err = format.WriteFormat(f, &te.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadTyp(t *typ.External) (err error) {
	format := typ.MakeFormatText(s.storeObjekten)

	var f *os.File

	if f, err = files.Open(t.FD.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	if _, err = format.ReadFormat(f, &t.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
