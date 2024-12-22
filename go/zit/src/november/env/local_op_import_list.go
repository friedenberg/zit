package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (u *Local) ImportList(
	list *sku.List,
	importer store.Importer,
) (err error) {
	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

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

			if err = coPrinter(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			continue
		}
	}

	if err = u.Unlock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if hasConflicts {
		err = store.ErrNeedsMerge
	}

	return
}
