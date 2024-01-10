package store_objekten

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/objekte_mode"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

func (s *Store) Import(sk *sku.Transacted) (co *sku.CheckedOut, err error) {
	co = sku.GetCheckedOutPool().Get()
	co.IsImport = true

	if err = co.External.Transacted.SetFromSkuLike(sk); err != nil {
		panic(err)
	}

	if err = sk.CalculateObjekteSha(); err != nil {
		co.SetError(err)
		err = nil
		return
	}

	_, err = s.GetBestandsaufnahmeStore().ReadOneEnnui(&sk.Metadatei.Sha)

	if err == nil {
		co.SetError(collections.ErrExists)
		return
	} else if errors.Is(err, objekte_store.ErrNotFoundEmpty) {
		err = nil
	} else {
		err = errors.Wrap(err)
		return
	}

	if err = s.ReadOneInto(sk.GetKennung(), &co.Internal); err != nil {
		if objekte_store.IsNotFound(err) {
			_, err = s.createOrUpdate(
				sk,
				sk.GetKennung(),
				nil,
				objekte_mode.ModeAddToBestandsaufnahme,
			)
		}

		err = errors.Wrap(err)
		return
	}

	if !co.Internal.Metadatei.Sha.IsNull() &&
		!co.Internal.Metadatei.Sha.Equals(&sk.Metadatei.Mutter) &&
		!co.Internal.Metadatei.Sha.Equals(&sk.Metadatei.Sha) {
		if err = s.importDoMerge(co); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = errors.Wrap(objekte_store.ErrLockRequired{
			Operation: fmt.Sprintf(
				"import %s",
				sk.GetGattung(),
			),
		})

		return
	}

	_, err = s.createOrUpdate(
		sk,
		sk.GetKennung(),
		&co.Internal,
		objekte_mode.ModeAddToBestandsaufnahme,
	)

	if err == collections.ErrExists {
		co.SetError(err)
		err = nil
	}

	return
}

var ErrNeedsMerge = errors.New("needs merge")

func (s *Store) importDoMerge(co *sku.CheckedOut) (err error) {
	co.SetError(ErrNeedsMerge)
	return
}
