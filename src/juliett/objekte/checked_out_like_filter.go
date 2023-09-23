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
) schnittstellen.FuncIter[*CheckedOut2] {
	if ms == nil {
		return collections.MakeWriterNoop[*CheckedOut2]()
	}

	return func(col *CheckedOut2) (err error) {
		internal := col.GetInternalLike()
		external := col.GetExternalLikePtr()

		g := gattung.Must(internal.GetSkuLike().GetGattung())

		var matcher matcher.Matcher
		ok := false

		if matcher, ok = ms.Get(g); !ok {
			err = iter.MakeErrStopIteration()
			return
		}

		if !matcher.ContainsMatchable(external) {
			err = iter.MakeErrStopIteration()
			return
		}

		return
	}
}
