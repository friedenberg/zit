package matcher

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func MakeFilterFromQuery(
	ms Query,
) schnittstellen.FuncIter[*sku.CheckedOut] {
	if ms == nil {
		return collections.MakeWriterNoop[*sku.CheckedOut]()
	}

	return func(col *sku.CheckedOut) (err error) {
		g := gattung.Must(col.Internal.GetSkuLike().GetGattung())

		var matcher Matcher
		ok := false

		if matcher, ok = ms.Get(g); !ok {
			err = iter.MakeErrStopIteration()
			return
		}

		if !matcher.ContainsMatchable(&col.External) {
			err = iter.MakeErrStopIteration()
			return
		}

		return
	}
}
