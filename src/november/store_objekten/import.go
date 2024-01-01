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
	if err = sk.CalculateObjekteSha(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var schwanz *sku.Transacted

	if schwanz, err = s.ReadOne(sk.GetKennung()); err != nil {
		if objekte_store.IsNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	if schwanz == nil {
		_, err = s.createOrUpdate(
			sk,
			sk.GetKennung(),
			nil,
			objekte_update_type.ModeAddToBestandsaufnahme,
		)

		return
	}

	err = s.GetVerzeichnisse().ExistsOneSha(&sk.Metadatei.Sha)

	if err == nil {
		return
	} else if errors.Is(err, objekte_store.ErrNotFoundEmpty) {
		err = nil
	} else {
		err = errors.Wrap(err)
		return
	}

	if !schwanz.Metadatei.Sha.Equals(&sk.Metadatei.Mutter) {
		if err = s.importDoMerge(sk, schwanz); err != nil {
			err = errors.Wrap(err)
			return
		}
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

	_, err = s.createOrUpdate(
		sk,
		sk.GetKennung(),
		schwanz,
		objekte_update_type.ModeAddToBestandsaufnahme,
	)

	return
}

func (s *Store) importDoMerge(
	sk, mutter *sku.Transacted,
) (err error) {
	return errors.Errorf("conflict: LOCAL: %s, REMOTE: %s", mutter, sk)
}
