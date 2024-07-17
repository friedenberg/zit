package chrome

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadCheckedOutFromItem(
	o sku.CommitOptions,
	k interfaces.ObjectId,
	em *item,
) (co *CheckedOut, err error) {
	co = GetCheckedOutPool().Get()

	if err = s.externalStoreInfo.FuncReadOneInto(k, &co.Internal); err != nil {
		// if collections.IsErrNotFound(err) {
		// 	// TODO mark status as new
		// 	err = nil
		// 	co.Internal.Kennung.ResetWith(&em.Kennung)
		// } else {
		// 	err = errors.Wrap(err)
		// 	return
		// }
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
		if err = co.Internal.SetFromSkuLike(sk); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	sku.DetermineState(co, false)

	return
}
