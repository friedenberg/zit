package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
)

func MakeFilterFromMetaSet(
	ms kennung.MetaSet,
) schnittstellen.FuncIter[CheckedOutLike] {
	if ms == nil {
		return collections.MakeWriterNoop[CheckedOutLike]()
	}

	return func(col CheckedOutLike) (err error) {
		internal := col.GetInternalLike()
		external := col.GetExternalLike()

		g := gattung.Must(internal.GetSkuLike().GetGattung())

		var matcher kennung.Matcher
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