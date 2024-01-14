package objekte_store

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/hotel/sku"
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
