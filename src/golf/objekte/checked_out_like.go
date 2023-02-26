package objekte

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type CheckedOutLike interface {
	GetInternal() TransactedLike // TODO-P0 rename to GetInternalLike
	GetExternal() ExternalLike   // TODO-P0 rename to GetExternalLike
	GetState() CheckedOutState
}

func MakeFilterFromMetaSet(
	ms kennung.MetaSet,
) schnittstellen.FuncIter[CheckedOutLike] {
	if ms == nil {
		return collections.MakeWriterNoop[CheckedOutLike]()
	}

	fds := ms.GetFDs()

	return func(col CheckedOutLike) (err error) {
		e := col.GetExternal()

		if fds.Contains(e.GetObjekteFD()) || fds.Contains(e.GetAkteFD()) {
			return
		}

		g := gattung.Must(col.GetInternal().GetDataIdentity().GetGattung())

		var ids kennung.Set
		ok := false

		if ids, ok = ms.Get(g); !ok {
			err = iter.MakeErrStopIteration()
			return
		}

		if ids.Sigil.IncludesCwd() && ids.Len() == 0 {
			return
		}

		id := col.GetInternal().GetSkuLike().GetId()

		if ids.Contains(id) {
			return
		}

		err = iter.MakeErrStopIteration()

		return
	}
}
