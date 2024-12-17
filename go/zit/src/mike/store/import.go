package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/store_fs"
)

var ErrNeedsMerge = errors.NewNormal("needs merge")

type BlobCopyResult struct {
	*sku.Transacted
	N int64
}

type Importer struct {
	*Store
	RemoteBlobStore    dir_layout.BlobStore
	BlobCopierDelegate interfaces.FuncIter[BlobCopyResult]
	ErrPrinter         interfaces.FuncIter[*sku.CheckedOut]
}

func (s Importer) Import(
	external *sku.Transacted,
) (co *sku.CheckedOut, err error) {
	if err = s.importBlobIfNecessary(external); err != nil {
		err = errors.Wrap(err)
		return
	}

	if external.GetGenre() == genres.InventoryList {
		if co, err = s.importInventoryList(external); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if co, err = s.importLeafSku(external); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s Importer) importInventoryList(
	el *sku.Transacted,
) (co *sku.CheckedOut, err error) {
	if el.GetGenre() == genres.InventoryList {
		if err = s.GetBlobStore().GetInventoryList().StreamInventoryListBlobSkus(
			el,
			func(sk *sku.Transacted) (err error) {
				if _, err = s.Import(
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
	}

	if co, err = s.importLeafSku(
		el,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return

	return
}

func (s Importer) importLeafSku(
	external *sku.Transacted,
) (co *sku.CheckedOut, err error) {
	co = store_fs.GetCheckedOutPool().Get()

	sku.Resetter.ResetWith(co.GetSkuExternal(), external)

	// if err = external.CalculateObjectShas(); err != nil {
	// 	co.SetError(err)
	// 	err = nil
	// 	return
	// }

	_, err = s.GetStreamIndex().ReadOneObjectIdTai(
		external.GetObjectId(),
		external.GetTai(),
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
	if err = s.ReadOneInto(external.GetObjectId(), co.GetSku()); err != nil {
		if collections.IsErrNotFound(err) {
			if err = s.tryRealizeAndOrStore(
				external,
				sku.CommitOptions{
					Clock:              co.GetSkuExternal(),
					Mode:               object_mode.ModeCommit,
					DontAddMissingTags: true,
					DontAddMissingType: true,
					ChangeIsHistorical: true,
				},
			); err != nil {
				err = errors.WrapExcept(err, collections.ErrExists)
			}
		} else {
			err = errors.Wrapf(err, "ObjectId: %s", external.GetObjectId())
		}

		return
	}

	if co.GetSku().Metadata.Sha().IsNull() {
		err = errors.Errorf("empty sha")
		return
	}

	var commitOptions sku.CommitOptions

	if commitOptions, err = s.MergeCheckedOutIfNecessary(co); err != nil {
		err = errors.Wrap(err)
		return
	}

	commitOptions.ChangeIsHistorical = true

	if err = s.tryRealizeAndOrStore(
		external,
		commitOptions,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Importer) importBlobIfNecessary(
	sk *sku.Transacted,
) (err error) {
	blobSha := sk.GetBlobSha()

	var n int64

	if n, err = dir_layout.CopyBlobIfNecessary(
		c.GetDirectoryLayout(),
		c.RemoteBlobStore,
		blobSha,
	); err != nil {
		if errors.Is(err, &dir_layout.ErrAlreadyExists{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if c.BlobCopierDelegate != nil {
		if err = c.BlobCopierDelegate(
			BlobCopyResult{Transacted: sk, N: n},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
