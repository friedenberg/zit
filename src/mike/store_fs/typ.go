package store_fs

import (
	"fmt"
	"os"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
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
			FD: kennung.FD{
				Path: fmt.Sprintf("%s.%s", t.Kennung(), s.erworben.FileExtensions.Typ),
			},
			// TODO-P2 move to central place
			Objekte: t.Objekte,
			Sku: sku.External[kennung.Typ, *kennung.Typ]{
				ObjekteSha: sha.Make(t.GetObjekteSha()),
				Kennung:    t.Sku.Kennung,
			},
		},
	}

	var f *os.File

	if f, err = files.CreateExclusiveWriteOnly(te.External.FD.Path); err != nil {
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

func (s *Store) ReadTyp(t *typ.External) (err error) {
	format := typ.MakeFormatText(s.storeObjekten)

	ops := objekte_store.MakeParseSaver[
		typ.Objekte,
		*typ.Objekte,
	](
		s.storeObjekten,
		s.storeObjekten,
		format,
	)

	var objekteSha schnittstellen.Sha

	if t.Objekte, objekteSha, err = ops.ParseAndSaveAkteAndObjekte(
		t.FD.Path,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.Sku.ObjekteSha = objekteSha.(sha.Sha)

	return
}
