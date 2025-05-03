package inventory_list_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	pkg_importer "code.linenisgreat.com/zit/go/zit/src/mike/importer"
)

func (store *Store) MakeImporter(
	options sku.ImporterOptions,
	storeOptions sku.StoreOptions,
) sku.Importer {
	importer := pkg_importer.Make(
		options,
		storeOptions,
		store.envRepo,
		store.getTypedBlobStore(),
		nil,
		nil,
		store,
	)

	return importer
}

func (store *Store) ImportList(
	list *sku.List,
	importer sku.Importer,
) (err error) {
	var hasConflicts bool

	oldPrinter := importer.GetCheckedOutPrinter()

	importer.SetCheckedOutPrinter(
		func(co *sku.CheckedOut) (err error) {
			if co.GetState() == checked_out_state.Conflicted {
				hasConflicts = true
			}

			return oldPrinter(co)
		},
	)

	for sk := range list.All() {
		if _, err = importer.Import(
			sk,
		); err != nil {
			if errors.Is(err, collections.ErrExists) {
				err = nil
			} else {
				err = errors.Wrapf(err, "Sku: %s", sku.String(sk))
				return
			}
		}
	}

	if hasConflicts {
		err = pkg_importer.ErrNeedsMerge
	}

	return
}
