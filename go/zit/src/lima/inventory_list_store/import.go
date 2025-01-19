package inventory_list_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/mike/importer"
)

func (store *Store) MakeImporter(
	options sku.ImporterOptions,
	storeOptions sku.StoreOptions,
) sku.Importer {
	importer := importer.Make(
		options,
		storeOptions,
		store.envRepo,
		store.typedBlobStore,
		nil,
		nil,
		store,
	)

	return importer
}

func (store *Store) ImportList(
	list *sku.List,
	i sku.Importer,
) (err error) {
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

	if hasConflicts {
		err = importer.ErrNeedsMerge
	}

	return
}
