package store

import (
	"fmt"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/objekte"
)

func (s *Store) Import(sk *sku.Transacted) (co *sku.CheckedOut, err error) {
	co = sku.GetCheckedOutPool().Get()
	co.IsImport = true

	if err = co.External.Transacted.SetFromSkuLike(sk); err != nil {
		panic(err)
	}

	if err = sk.CalculateObjekteShas(); err != nil {
		co.SetError(err)
		err = nil
		return
	}

	_, err = s.GetVerzeichnisse().ReadOneEnnui(sk.Metadatei.Sha())

	if err == nil {
		co.SetError(collections.ErrExists)
		return
	} else if collections.IsErrNotFound(err) {
		err = nil
	} else {
		err = errors.Wrap(err)
		return
	}

	if err = s.ReadOneInto(sk.GetKennung(), &co.Internal); err != nil {
		if collections.IsErrNotFound(err) {
			_, err = s.createOrUpdate(
				sk,
				sk.GetKennung(),
				objekte_mode.ModeAddToBestandsaufnahme,
			)
		}

		err = errors.Wrap(err)
		return
	}

	if co.Internal.Metadatei.Sha().IsNull() {
		err = errors.Errorf("empty sha")
		return
	}

	if !co.Internal.Metadatei.Sha().IsNull() &&
		!co.Internal.Metadatei.Sha().Equals(sk.Metadatei.Mutter()) &&
		!co.Internal.Metadatei.Sha().Equals(sk.Metadatei.Sha()) {
		if err = s.importDoMerge(co); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = errors.Wrap(objekte.ErrLockRequired{
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
