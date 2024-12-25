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
	sku.ParentNegotiator
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
}

func (importer Importer) importLeafSku(
	external *sku.Transacted,
) (co *sku.CheckedOut, err error) {
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
	if err = importer.ReadOneInto(co.GetSkuExternal().GetObjectId(), co.GetSku()); err != nil {
		if collections.IsErrNotFound(err) {
			if err = importer.tryRealizeAndOrStore(
				co.GetSkuExternal(),
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

	var commitOptions sku.CommitOptions

	if commitOptions, err = importer.MergeCheckedOutIfNecessary(
		co,
		importer.ParentNegotiator,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	commitOptions.ChangeIsHistorical = true

	if err = importer.tryRealizeAndOrStore(
		co.GetSkuExternal(),
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
	if c.RemoteBlobStore == nil {
		err = errors.Errorf("nil blob store")
		return
	}

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
