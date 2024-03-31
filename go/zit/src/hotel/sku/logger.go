package sku

import (
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/collections"
)

type Logger struct {
	New, Updated, Unchanged, Archived schnittstellen.FuncIter[*Transacted]
}

func (l Logger) NewOrUpdated(
	err error,
) schnittstellen.FuncIter[*Transacted] {
	if collections.IsErrNotFound(err) {
		return l.New
	} else {
		return l.Updated
	}
}
