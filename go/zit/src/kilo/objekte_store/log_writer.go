package objekte_store

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/hotel/sku"
)

type LogWriter struct {
	New, Updated, Unchanged, Archived schnittstellen.FuncIter[*sku.Transacted]
}

func (l LogWriter) NewOrUpdated(
	err error,
) schnittstellen.FuncIter[*sku.Transacted] {
	if collections.IsErrNotFound(err) {
		return l.New
	} else {
		return l.Updated
	}
}
