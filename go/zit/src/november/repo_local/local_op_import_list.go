package repo_local

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (u *Repo) ImportList(
	list *sku.List,
	importer store.Importer,
) (err error) {
	u.Must(u.Lock)

	coPrinter := u.PrinterCheckedOut(box_format.CheckedOutHeaderState{})

	var co *sku.CheckedOut
	hasConflicts := false

	for {
		sk, ok := list.Pop()

		if !ok {
			break
		}

		if co, err = importer.Import(
			sk,
		); err != nil {
			err = errors.Wrapf(err, "Sku: %s", sk)
			return
		}

		if co.GetState() == checked_out_state.Conflicted {
			hasConflicts = true

			if !importer.DontPrint {
				if err = coPrinter(co); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			continue
		}
	}

	u.Must(u.Unlock)

	if hasConflicts {
		err = store.ErrNeedsMerge
	}

	return
}
