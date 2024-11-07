package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) ToSliceFilesZettelen(
	cos sku.CheckedOutLikeSet,
) (out []string, err error) {
	return quiter.DerivedValues(
		cos,
		func(col sku.CheckedOutLike) (e string, err error) {
			var fds *sku.FSItem

			if fds, err = s.ReadFSItemFromExternal(col.GetSkuExternalLike()); err != nil {
				err = errors.Wrap(err)
				return
			}

			e = fds.Object.GetPath()

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}

func (s *Store) ToSliceFilesBlobs(
	cos sku.CheckedOutLikeSet,
) (out []string, err error) {
	return quiter.DerivedValues(
		cos,
		func(col sku.CheckedOutLike) (e string, err error) {
			var fds *sku.FSItem

			if fds, err = s.ReadFSItemFromExternal(col.GetSkuExternalLike()); err != nil {
				err = errors.Wrap(err)
				return
			}

			e = fds.Blob.GetPath()

			if e == "" {
				err = collections.MakeErrStopIteration()
				return
			}

			return
		},
	)
}
