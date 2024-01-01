package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

func (s *Store) Import(
	sk *sku.Transacted,
) (err error) {
  // log.Debug().Printf("%s", sk.StringKennungTai())
	if err = sk.CalculateObjekteSha(); err != nil {
		err = errors.Wrap(err)
		return
	}

	err = s.GetVerzeichnisse().ExistsOneSha(&sk.Metadatei.Sha)

	if errors.Is(err, objekte_store.ErrNotFoundEmpty) {
		err = nil
	} else {
		err = errors.Wrap(err)
		return
	}

	if _, err = s.CreateOrUpdate(sk, sk.GetKennung()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
