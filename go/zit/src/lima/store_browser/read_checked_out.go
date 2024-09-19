package store_browser

import (
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ReadIntoCheckedOutFromTransacted(
	sk *sku.Transacted,
	co *sku.CheckedOut,
) (err error) {
	if &co.Internal != sk {
		sku.Resetter.ResetWith(&co.Internal, sk)
	}

	sku.DetermineState(co, false)

	return
}
