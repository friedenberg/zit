package store

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
)

func (s *Store) Import(external *sku.Transacted) (co *store_fs.CheckedOut, err error) {
	co = store_fs.GetCheckedOutPool().Get()
	co.IsImport = true

	if err = co.External.Transacted.SetFromSkuLike(external); err != nil {
		panic(err)
	}

	if err = external.CalculateObjectShas(); err != nil {
		co.SetError(err)
		err = nil
		return
	}

	_, err = s.GetStreamIndex().ReadOneObjectIdTai(
		external.GetObjectId(),
		external.GetTai(),
	)

	if err == nil {
		co.SetError(collections.ErrExists)
		return
	} else if collections.IsErrNotFound(err) {
		err = nil
	} else {
		err = errors.Wrap(err)
		return
	}

	if err = s.ReadOneInto(external.GetObjectId(), &co.Internal); err != nil {
		if collections.IsErrNotFound(err) {
			err = s.tryRealizeAndOrStore(
				external,
				ObjekteOptions{
					Clock: &co.External.Transacted,
					Mode:  objekte_mode.ModeCommit,
				},
			)

			err = errors.WrapExcept(err, collections.ErrExists)
		}

		err = errors.Wrap(err)

		return
	}

	if co.Internal.Metadata.Sha().IsNull() {
		err = errors.Errorf("empty sha")
		return
	}

	if !co.Internal.Metadata.Sha().IsNull() &&
		!co.Internal.Metadata.Sha().Equals(external.Metadata.Mutter()) &&
		!co.Internal.Metadata.Sha().Equals(external.Metadata.Sha()) {
		if err = s.importDoMerge(co); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = errors.Wrap(file_lock.ErrLockRequired{
			Operation: fmt.Sprintf(
				"import %s",
				external.GetGenre(),
			),
		})

		return
	}

	if err = s.tryRealizeAndOrStore(
		external,
		ObjekteOptions{
			Mode: objekte_mode.ModeCommit,
		},
	); err == collections.ErrExists {
		co.SetError(err)
		err = nil
	}

	return
}

var ErrNeedsMerge = errors.New("needs merge")

func (s *Store) importDoMerge(co *store_fs.CheckedOut) (err error) {
	co.SetError(ErrNeedsMerge)
	return
}
