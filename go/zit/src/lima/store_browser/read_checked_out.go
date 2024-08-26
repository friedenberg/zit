package store_browser

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadCheckedOutFromItem(
	o sku.CommitOptions,
	k interfaces.ObjectId,
	em *browserItemRaw,
) (co *CheckedOut, err error) {
	co = GetCheckedOutPool().Get()

	if err = s.externalStoreInfo.FuncReadOneInto(
		k.String(),
		&co.Internal,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.ReadIntoCheckedOutFromTransacted(&co.Internal, co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) ReadIntoCheckedOutFromTransacted(
	sk *sku.Transacted,
	co *CheckedOut,
) (err error) {
	if &co.Internal != sk {
		sku.Resetter.ResetWith(&co.Internal, sk)
	}

	sku.DetermineState(co, false)

	return
}
