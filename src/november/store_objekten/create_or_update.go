package store_objekten

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/checkout_mode"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/objekte_mode"
	"github.com/friedenberg/zit/src/charlie/checkout_options"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/to_merge"
	"github.com/friedenberg/zit/src/kilo/objekte_store"
)

type CreateOrUpdateDelegate struct {
	New       schnittstellen.FuncIter[*sku.Transacted]
	Updated   schnittstellen.FuncIter[*sku.Transacted]
	Unchanged schnittstellen.FuncIter[*sku.Transacted]
}

func (s *Store) CreateOrUpdateCheckedOut(
	co *sku.CheckedOut,
) (transactedPtr *sku.Transacted, err error) {
	kennungPtr := &co.External.Kennung

	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: fmt.Sprintf("create or update %s", kennungPtr),
		}

		return
	}

	transactedPtr = sku.GetTransactedPool().Get()

	if err = transactedPtr.SetFromSkuLike(&co.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = transactedPtr.SetAkteSha(co.External.GetAkteSha()); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P2: determine why Metadatei.Etiketten can be nil
	if transactedPtr.Metadatei.EqualsSansTai(&co.Internal.Metadatei) {
		transactedPtr = &co.Internal

		if err = s.handleUnchanged(transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.ReadOne(kennungPtr); err != nil {
		if objekte_store.IsNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	} else {
		mu := &mutter.Metadatei.Sha
		if err = transactedPtr.Metadatei.Mutter.SetShaLike(
			mu,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.handleUpdated(
		transactedPtr,
		objekte_mode.ModeCommit,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CreateOrUpdateTransacted(
	in *sku.Transacted,
) (out *sku.Transacted, err error) {
	return s.CreateOrUpdate(in, in.GetKennung())
}

func (s *Store) createOrUpdate(
	mg metadatei.Getter,
	kennungPtr kennung.Kennung,
	mutter *sku.Transacted,
	updateType objekte_mode.Mode,
) (transactedPtr *sku.Transacted, err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var m *metadatei.Metadatei

	if mg != nil {
		m = mg.GetMetadatei()
	} else {
		m = metadatei.GetPool().Get()
		defer metadatei.GetPool().Put(m)
	}

	transactedPtr = sku.GetTransactedPool().Get()
	metadatei.Resetter.ResetWith(&transactedPtr.Metadatei, m)

	if err = transactedPtr.Kennung.SetWithKennung(kennungPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter != nil {
		transactedPtr.Kopf = mutter.GetKopf()
		mu := &mutter.Metadatei.Sha

		if err = transactedPtr.Metadatei.Mutter.SetShaLike(
			mu,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		errors.TodoP4("determine if this is necessary any more")
		// transactedPtr.Sku.Kopf = s.common.GetTransaktion().Time
	}

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter != nil &&
		kennung.Equals(transactedPtr.GetKennung(), mutter.GetKennung()) &&
		transactedPtr.Metadatei.EqualsSansTai(&mutter.Metadatei) {
		if err = transactedPtr.SetFromSkuLike(mutter); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = s.handleUnchanged(transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = s.readExternalAndMergeIfNecessary(transactedPtr, mutter); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.handleNewOrUpdateWithMutter(
		transactedPtr,
		mutter,
		updateType,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CreateOrUpdate(
	mg metadatei.Getter,
	kennungPtr kennung.Kennung,
) (transactedPtr *sku.Transacted, err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.ReadOne(kennungPtr); err != nil {
		if objekte_store.IsNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return s.createOrUpdate(
		mg,
		kennungPtr,
		mutter,
		objekte_mode.ModeCommit,
	)
}

func (s *Store) readExternalAndMergeIfNecessary(
	transactedPtr, mutter *sku.Transacted,
) (err error) {
	if mutter == nil {
		return
	}

	var co *sku.CheckedOut

	if co, err = s.ReadOneExternalFS(mutter); err != nil {
		err = nil
		return
	}

	defer sku.GetCheckedOutPool().Put(co)

	mutterEqualsExternal := co.InternalAndExternalEqualsSansTai()

	var mode checkout_mode.Mode

	if mode, err = co.External.GetCheckoutMode(); err != nil {
		err = errors.Wrap(err)
		return
	}

	op := checkout_options.Options{
		CheckoutMode: mode,
		Force:        true,
	}

	if mutterEqualsExternal {
		if co, err = s.CheckoutOne(op, transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

		sku.GetCheckedOutPool().Put(co)

		return
	}

	transactedPtrCopy := sku.GetTransactedPool().Get()
	defer sku.GetTransactedPool().Put(transactedPtrCopy)

	if err = transactedPtrCopy.SetFromSkuLike(transactedPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	tm := to_merge.Sku{
		Left:   transactedPtrCopy,
		Middle: &co.Internal,
		Right:  &co.External.Transacted,
	}

	var merged sku.ExternalFDs

	merged, err = s.merge(tm)

	switch {
	case errors.Is(err, &to_merge.ErrMergeConflict{}):
		if err = tm.WriteConflictMarker(
			s.GetStandort(),
			s.GetKonfig().GetStoreVersion(),
			s.GetObjekteFormatOptions(),
			co.External.FDs.MakeConflictMarker(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

	case err != nil:
		err = errors.Wrap(err)
		return

	default:
		src := merged.Objekte.GetPath()
		dst := co.External.FDs.Objekte.GetPath()

		if err = files.Rename(src, dst); err != nil {
			return
		}
	}

	return
}

func (s *Store) CreateOrUpdateAkteSha(
	kennungPtr kennung.Kennung,
	sh schnittstellen.ShaLike,
) (transactedPtr *sku.Transacted, err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = objekte_store.ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.ReadOne(kennungPtr); err != nil {
		if objekte_store.IsNotFound(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	transactedPtr = sku.GetTransactedPool().Get()

	if mutter == nil {
		if err = transactedPtr.Kennung.SetWithKennung(kennungPtr); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		sku.TransactedResetter.ResetWith(transactedPtr, mutter)
	}

	transactedPtr.SetAkteSha(sh)

	return s.createOrUpdate(
		transactedPtr,
		kennungPtr,
		mutter,
		objekte_mode.ModeCommit,
	)
}

func (s *Store) handleNewOrUpdateWithMutter(
	sk, mutter *sku.Transacted,
	updateType objekte_mode.Mode,
) (err error) {
	if err = iter.Chain(
		sk,
		func(t1 *sku.Transacted) error {
			if mutter == nil {
				return s.handleNew(t1, updateType)
			} else {
				return s.handleUpdated(t1, updateType)
			}
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
