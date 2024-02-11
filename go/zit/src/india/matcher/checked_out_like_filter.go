package matcher

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

func MakeFilterFromQuery(
	ms Query,
) schnittstellen.FuncIter[*sku.CheckedOut] {
	if ms == nil {
		return collections.MakeWriterNoop[*sku.CheckedOut]()
	}

	return func(col *sku.CheckedOut) (err error) {
		if !ms.ContainsMatchable(&col.External.Transacted) {
			err = iter.MakeErrStopIteration()
			return
		}

		return
	}
}
