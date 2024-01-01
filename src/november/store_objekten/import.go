package store_objekten

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/objekte_update_type"
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

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: fmt.Sprintf(
				"import %s",
				sk.GetGattung(),
			),
		}

		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.ReadOne(sk.GetKennung()); err != nil {
		if objekte_store.IsNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	_, err = s.createOrUpdate(
		sk,
		sk.GetKennung(),
		mutter,
		objekte_update_type.ModeAddToBestandsaufnahme,
	)

	return
}
