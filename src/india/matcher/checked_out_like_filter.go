package matcher

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/hotel/sku"
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
