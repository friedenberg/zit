package store_fs

import (
	"fmt"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/hotel/typ"
)

func (s *Store) CheckinTyp(p string) (t *typ.Transacted, err error) {
	return
}

func (s *Store) WriteTyp(t *typ.Transacted) (te *typ.CheckedOut, err error) {
	te = &typ.CheckedOut{
		Internal: *t,
		External: typ.External{
			Sku: sku.External[kennung.Typ, *kennung.Typ]{
				ObjekteSha: sha.Make(t.GetObjekteSha()),
				Kennung:    t.Sku.Kennung,
				ObjekteFD: kennung.FD{
					Path: fmt.Sprintf("%s.%s", t.Kennung(), s.erworben.FileExtensions.Typ),
				},
			},
			// TODO-P2 move to central place
			Objekte: t.Objekte,
		},
	}

	var f *os.File

	if f, err = files.CreateExclusiveWriteOnly(
		te.External.GetObjekteFD().Path,
	); err != nil {
		if errors.IsExist(err) {
			err = s.ReadTyp(&te.External)
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	defer errors.Deferred(&err, f.Close)

	format := typ.MakeFormatText(s.storeObjekten)

	if _, err = format.WriteFormat(f, &te.External.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadTypFromFile(p string) (t typ.External, err error) {
	format := typ.MakeFormatText(s.storeObjekten)

	ops := objekte_store.MakeParseSaver[
		typ.Objekte,
		*typ.Objekte,
		kennung.Typ,
		*kennung.Typ,
	](
		s.storeObjekten,
		s.storeObjekten,
		format,
	)

	if t.Objekte, t.Sku, err = ops.ParseAndSaveAkteAndObjekte(
		p,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadTyp(t *typ.External) (err error) {
	*t, err = s.ReadTypFromFile(t.GetObjekteFD().Path)
	return
}
