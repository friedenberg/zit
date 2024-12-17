package env

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (u *Local) ImportListFromRemoteBlobStore(
	list *sku.List,
	remoteBlobStore dir_layout.BlobStore,
	printCopies bool,
) (err error) {
	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	coPrinter := u.PrinterCheckedOut(box_format.CheckedOutHeaderState{})

	importer := store.Importer{
		Store:           u.GetStore(),
		ErrPrinter:      coPrinter,
		RemoteBlobStore: remoteBlobStore,
	}

	if printCopies {
		importer.BlobCopierDelegate = func(result store.BlobCopyResult) error {
			// TODO switch to Err and fix test
			return ui.Out().Printf(
				"copied Blob %s (%d bytes)",
				result.GetBlobSha(),
				result.N,
			)
		}
	}

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
