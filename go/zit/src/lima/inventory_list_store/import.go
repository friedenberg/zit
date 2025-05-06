package inventory_list_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
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

	checkedOutPrinter := importer.GetCheckedOutPrinter()

	importer.SetCheckedOutPrinter(
		func(co *sku.CheckedOut) (err error) {
			if co.GetState() == checked_out_state.Conflicted {
				hasConflicts = true
			}

			return checkedOutPrinter(co)
		},
	)

	importErrors := errors.MakeMulti()
	missingBlobs := sku.MakeListCheckedOut()

	for sk := range list.All() {
		checkedOut, importError := importer.Import(sk)

		func() {
			defer sku.GetCheckedOutPool().Put(checkedOut)

			if importError == nil {
				// TODO eliminate condition where checkedOut can be nil and err can be
				// nil
				if checkedOut != nil &&
					checkedOut.GetState() == checked_out_state.Conflicted {
					hasConflicts = true
				}

				return
			}

			if errors.Is(importError, collections.ErrExists) {
				return
			}

			if genres.IsErrUnsupportedGenre(importError) {
				return
			}

			if env_dir.IsErrBlobMissing(importError) {
				checkedOut := sku.GetCheckedOutPool().Get()
				sku.TransactedResetter.ResetWith(checkedOut.GetSkuExternal(), sk)
				checkedOut.SetState(checked_out_state.Untracked)

				missingBlobs.Add(checkedOut)

				return
			}

			importErrors.Add(errors.Wrapf(err, "Sku: %s", sku.String(sk)))
		}()
	}

	checkedOutPrinter = store.ui.CheckedOutCheckedOut

	if missingBlobs.Len() > 0 {
		ui.Err().Printf(
			"could not import the %d objects (blobs missing):\n",
			missingBlobs.Len(),
		)

		for missing := range missingBlobs.All() {
			if err = checkedOutPrinter(missing); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	if hasConflicts {
		importErrors.Add(pkg_importer.ErrNeedsMerge)
	}

	if importErrors.Len() > 0 {
		err = importErrors
	}

	return
}
