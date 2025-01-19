package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/mike/importer"
)

func (u *Repo) ImportList(
	list *sku.List,
	i importer.Importer,
) (err error) {
	u.Must(u.Lock)

	var hasConflicts bool

	oldPrinter := i.GetCheckedOutPrinter()

	i.SetCheckedOutPrinter(
		func(co *sku.CheckedOut) (err error) {
			if co.GetState() == checked_out_state.Conflicted {
				hasConflicts = true
			}

			return oldPrinter(co)
		},
	)

	for {
		sk, ok := list.Pop()

		if !ok {
			break
		}

		if _, err = i.Import(
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
		err = importer.ErrNeedsMerge
	}

	return
}
