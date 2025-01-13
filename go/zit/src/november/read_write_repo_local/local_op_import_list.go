package read_write_repo_local

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/india/sku"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (u *Repo) ImportList(
	list *sku.List,
	importer store.Importer,
) (err error) {
	u.Must(u.Lock)

	var hasConflicts bool

	oldPrinter := importer.CheckedOutPrinter

	importer.CheckedOutPrinter = func(co *sku.CheckedOut) (err error) {
		if co.GetState() == checked_out_state.Conflicted {
			hasConflicts = true
		}

		return oldPrinter(co)
	}

	for {
		sk, ok := list.Pop()

		if !ok {
			break
		}

		if _, err = importer.Import(
			sk,
		); err != nil {
			if errors.Is(err, collections.ErrExists) {
				err = nil
			} else {
				err = errors.Wrapf(err, "Sku: %s", sk)
				return
			}
		}
	}

	u.Must(u.Unlock)

	if hasConflicts {
		err = store.ErrNeedsMerge
	}

	return
}
