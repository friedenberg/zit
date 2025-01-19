package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
)

var ErrNeedsMerge = errors.NewNormal("needs merge")

type ImporterOptions struct {
	ExcludeObjects      bool
	RemoteBlobStore     interfaces.BlobStore
	PrintCopies         bool
	AllowMergeConflicts bool
	BlobCopierDelegate  interfaces.FuncIter[sku.BlobCopyResult]
	ParentNegotiator    sku.ParentNegotiator
	CheckedOutPrinter   interfaces.FuncIter[*sku.CheckedOut]
}

func (store *Store) MakeImporter(
	options ImporterOptions,
	storeOptions sku.StoreOptions,
) (importer Importer) {
	importer = Importer{
		Store:               store,
		ExcludeObjects:      options.ExcludeObjects,
		RemoteBlobStore:     options.RemoteBlobStore,
		BlobCopierDelegate:  options.BlobCopierDelegate,
		AllowMergeConflicts: options.AllowMergeConflicts,
		ParentNegotiator:    options.ParentNegotiator,
		CheckedOutPrinter:   options.CheckedOutPrinter,
		StoreOptions:        storeOptions,
	}

	if importer.BlobCopierDelegate == nil &&
		importer.RemoteBlobStore != nil &&
		options.PrintCopies {
		importer.BlobCopierDelegate = sku.MakeBlobCopierDelegate(
			store.envRepo.GetUI(),
		)
	}

	return
}

type Importer struct {
	*Store
	ExcludeObjects      bool
	RemoteBlobStore     interfaces.BlobStore
	BlobCopierDelegate  interfaces.FuncIter[sku.BlobCopyResult]
	StoreOptions        sku.StoreOptions
	AllowMergeConflicts bool
	sku.ParentNegotiator
	CheckedOutPrinter interfaces.FuncIter[*sku.CheckedOut]
}

func (importer Importer) Import(
	external *sku.Transacted,
) (co *sku.CheckedOut, err error) {
	if err = importer.ImportBlobIfNecessary(external); err != nil {
		err = errors.Wrap(err)
		return
	}

	if external.GetGenre() == genres.InventoryList {
		if co, err = importer.importInventoryList(external); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if co, err = importer.importLeafSku(external); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (importer Importer) importInventoryList(
	el *sku.Transacted,
) (co *sku.CheckedOut, err error) {
	if importer.RemoteBlobStore == nil {
		err = errors.Errorf("RemoteBlobStore is nil")
		return
	}

	if el.GetGenre() != genres.InventoryList {
		err = errors.Errorf(
			"Expected genre %q but got %q",
			genres.InventoryList,
			el.GetGenre(),
		)
		return
	}

	if err = importer.GetBlobStore().GetInventoryList().StreamInventoryListBlobSkus(
		el,
		func(sk *sku.Transacted) (err error) {
			if _, err = importer.Import(
				sk,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if co, err = importer.importLeafSku(
		el,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (importer Importer) importLeafSku(
	external *sku.Transacted,
) (co *sku.CheckedOut, err error) {
	if importer.ExcludeObjects {
		return
	}

	// TODO address this terrible hack? How should config objects be handled by
	// remotes?
	if external.GetGenre() == genres.Config {
		return
	}

	co = store_fs.GetCheckedOutPool().Get()

	sku.Resetter.ResetWith(co.GetSkuExternal(), external)

	if err = co.GetSkuExternal().CalculateObjectShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	_, err = importer.GetStreamIndex().ReadOneObjectIdTai(
		co.GetSkuExternal().GetObjectId(),
		co.GetSkuExternal().GetTai(),
	)

	if err == nil {
		err = collections.ErrExists
		return
	} else if collections.IsErrNotFound(err) {
		err = nil
	} else {
		err = errors.Wrap(err)
		return
	}

	ui.TodoP4("cleanup")
	if err = importer.ReadOneInto(
		co.GetSkuExternal().GetObjectId(),
		co.GetSku(),
	); err != nil {
		if collections.IsErrNotFound(err) {
			if err = importer.tryRealizeAndOrStore(
				co.GetSkuExternal(),
				sku.CommitOptions{
					Clock:              co.GetSkuExternal(),
					StoreOptions:       importer.StoreOptions,
					DontAddMissingTags: true,
					DontAddMissingType: true,
				},
			); err != nil {
				err = errors.WrapExcept(err, collections.ErrExists)
			}
		} else {
			err = errors.Wrapf(err, "ObjectId: %s", external.GetObjectId())
		}

		return
	}

	var commitOptions sku.CommitOptions

	// TODO extra commit option setting into its own function
	if commitOptions, err = importer.MergeCheckedOutIfNecessary(
		co,
		importer.ParentNegotiator,
		importer.AllowMergeConflicts,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if co.GetState() == checked_out_state.Conflicted {
		if !importer.AllowMergeConflicts {
			if err = importer.CheckedOutPrinter(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	commitOptions.Validate = false

	if err = importer.tryRealizeAndOrStore(
		co.GetSkuExternal(),
		commitOptions,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = importer.CheckedOutPrinter(co); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Importer) ImportBlobIfNecessary(
	sk *sku.Transacted,
) (err error) {
	blobSha := sk.GetBlobSha()

	if c.RemoteBlobStore == nil {
		// when this is a dumb HTTP remote, we expect local to push the missing
		// objects to us after the import call

		n := int64(-1)

		if c.GetDirectoryLayout().HasBlob(blobSha) {
			n = -2
		}

		if c.BlobCopierDelegate != nil {
			if err = c.BlobCopierDelegate(
				sku.BlobCopyResult{
					Transacted: sk,
					Sha:        blobSha,
					N:          n,
				},
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		return
	}

	var n int64

	if n, err = env_repo.CopyBlobIfNecessary(
		c.GetDirectoryLayout(),
		c.GetDirectoryLayout(),
		c.RemoteBlobStore,
		blobSha,
	); err != nil {
		if errors.Is(err, &env_dir.ErrAlreadyExists{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if c.BlobCopierDelegate != nil {
		if err = c.BlobCopierDelegate(
			sku.BlobCopyResult{
				Transacted: sk,
				Sha:        blobSha,
				N:          n,
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
