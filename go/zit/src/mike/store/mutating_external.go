package store

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/delta/file_lock"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
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

func (s *Store) CreateOrUpdateTransacted(
	in *sku.Transacted,
) (out *sku.Transacted, err error) {
	if in.Kennung.IsEmpty() {
		if in.GetGattung() != gattung.Zettel {
			err = errors.Errorf("only Zettel is supported")
			return
		}

		if out, err = s.Create(in); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if out, err = s.CreateOrUpdate(in, in.GetKennung()); err != nil {
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
		err = file_lock.ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	return s.createOrUpdate(
		mg,
		kennungPtr,
		objekte_mode.ModeCommit,
	)
}

func (s *Store) CreateOrUpdateAkteSha(
	kennungPtr kennung.Kennung,
	sh schnittstellen.ShaLike,
) (transactedPtr *sku.Transacted, err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var mutter *sku.Transacted

	if mutter, err = s.ReadOne(kennungPtr); err != nil {
		if collections.IsErrNotFound(err) {
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
		objekte_mode.ModeCommit,
	)
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

	if _, err = s.CreateOrUpdate(
		&mutter.Metadatei,
		sk.GetKennung(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

//   _____    _   _       _
//  |__  /___| |_| |_ ___| |
//    / // _ \ __| __/ _ \ |
//   / /|  __/ |_| ||  __/ |
//  /____\___|\__|\__\___|_|
//

func (s *Store) CreateWithAkteString(
	mg metadatei.Getter,
	akteString string,
) (tz *sku.Transacted, err error) {
	var aw sha.WriteCloser

	if aw, err = s.GetStandort().AkteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, err = io.WriteString(aw, akteString); err != nil {
		err = errors.Wrap(err)
		return
	}

	m := mg.GetMetadatei()
	m.SetAkteSha(aw)

	defer errors.DeferredCloser(&err, aw)

	if tz, err = s.Create(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) Create(
	mg metadatei.Getter,
) (tz *sku.Transacted, err error) {
	if !s.GetStandort().GetLockSmith().IsAcquired() {
		err = file_lock.ErrLockRequired{
			Operation: "create",
		}

		return
	}

	if mg.GetMetadatei().IsEmpty() {
		err = errors.Normalf("zettel is empty")
		return
	}

	if s.protoZettel.Equals(mg.GetMetadatei()) {
		err = errors.Normalf("zettel matches protozettel")
		return
	}

	m := mg.GetMetadatei()
	s.protoZettel.Apply(m)

	if err = s.GetKonfig().ApplyToNewMetadatei(
		m,
		s.GetAkten().GetTypV0(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	var ken *kennung.Hinweis

	if ken, err = s.kennungIndex.CreateHinweis(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if tz, err = s.makeSku(
		m,
		ken,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.tryCommit(
		tz,
		objekte_mode.ModeCommit,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
