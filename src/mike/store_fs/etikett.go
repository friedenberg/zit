package store_fs

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
)

func (s *Store) ReadEtikettFromFile(p string) (t etikett.External, err error) {
	format := etikett.MakeFormatText(s.storeObjekten)

	ops := objekte_store.MakeParseSaver[
		etikett.Objekte,
		*etikett.Objekte,
		kennung.Etikett,
		*kennung.Etikett,
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

func (s *Store) ReadEtikett(t *etikett.External) (err error) {
	*t, err = s.ReadEtikettFromFile(t.GetObjekteFD().Path)
	return
}
