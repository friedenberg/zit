package store

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO-P2 add support for quiet reindexing
func (s *Store) Reindex() (err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: "reindex",
		}

		return
	}

	if err = s.ResetIndexes(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetStandort().ResetVerzeichnisse(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetVerzeichnisse().Initialize(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetBestandsaufnahmeStore().ReadAllSkus(
		s.reindexOne,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CreateOrUpdateFromTransacted(
	in *sku.Transacted,
	mode objekte_mode.Mode,
) (err error) {
	mode.Add(
		objekte_mode.ModeCommit,
		objekte_mode.ModeApplyProto,
	)

	if err = s.tryRealizeAndOrStore(in, ObjekteOptions{Mode: mode}); err != nil {
		err = errors.WrapExcept(err, collections.ErrExists)
		return
	}

	return
}

func (s *Store) CreateOrUpdateAkteSha(
	k kennung.Kennung,
	sh interfaces.ShaLike,
) (t *sku.Transacted, err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				k.GetGenre(),
			),
		}

		return
	}

	t = sku.GetTransactedPool().Get()

	if err = t.Kennung.SetWithKennung(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.ReadOneInto(k, t); err != nil {
		if collections.IsErrNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	t.SetAkteSha(sh)

	if err = s.tryRealizeAndOrStore(
		t,
		ObjekteOptions{Mode: objekte_mode.ModeCommit},
	); err != nil {
		err = errors.WrapExcept(err, collections.ErrExists)
		return
	}

	return
}

func (s *Store) RevertTo(
	sk *sku.Transacted,
	sh *sha.Sha,
) (err error) {
	if sh.IsNull() {
		err = errors.Errorf("cannot revert to null")
		return
	}

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: "update many metadatei",
		}

		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.ReadOneEnnui(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	mutter.Metadatei.Mutter().ResetWith(sk.Metadatei.Sha())

	if err = mutter.CalculateObjekteShas(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer sku.GetTransactedPool().Put(mutter)

	if err = s.tryRealizeAndOrStore(
		mutter,
		ObjekteOptions{Mode: objekte_mode.ModeCommit},
	); err != nil {
		err = errors.WrapExcept(err, collections.ErrExists)
		return
	}

	return
}
