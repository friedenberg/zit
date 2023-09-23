package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/india/matcher"
)

func MakeFilterFromMetaSet(
	ms matcher.Query,
) schnittstellen.FuncIter[*CheckedOut] {
	if ms == nil {
		return collections.MakeWriterNoop[*CheckedOut]()
	}

	return func(col *CheckedOut) (err error) {
		g := gattung.Must(col.Internal.GetSkuLike().GetGattung())

		var matcher matcher.Matcher
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
