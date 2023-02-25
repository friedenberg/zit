package store_fs

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kilo/cwd"
)

func (s *Store) ReadEtikett(sem cwd.Etikett) (t etikett.External, err error) {
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
		sem,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
