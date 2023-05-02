package objekte_store

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
	"github.com/friedenberg/zit/src/india/konfig"
)

type CreateOrUpdateDelegate[T any] struct {
	New       schnittstellen.FuncIter[T]
	Updated   schnittstellen.FuncIter[T]
	Unchanged schnittstellen.FuncIter[T]
}

type createOrUpdate[
	T objekte.Akte[T],
	T1 objekte.AktePtr[T],
	T2 kennung.KennungLike[T2],
	T3 kennung.KennungLikePtr[T2],
	T4 any,
	T5 objekte.VerzeichnissePtr[T4, T],
] struct {
	clock                     ts.Clock
	ls                        schnittstellen.LockSmith
	of                        schnittstellen.ObjekteWriterFactory
	af                        schnittstellen.AkteWriterFactory
	reader                    TransactedReader[T3, *objekte.Transacted[T, T1, T2, T3, T4, T5]]
	delegate                  CreateOrUpdateDelegate[*objekte.Transacted[T, T1, T2, T3, T4, T5]]
	matchableAdder            kennung.MatchableAdder
	persistentMetadateiFormat persisted_metadatei_format.Format
	kg                        konfig.Getter
}

func MakeCreateOrUpdate[
	T objekte.Akte[T],
	T1 objekte.AktePtr[T],
	T2 kennung.KennungLike[T2],
	T3 kennung.KennungLikePtr[T2],
	T4 any,
	T5 objekte.VerzeichnissePtr[T4, T],
](
	clock ts.Clock,
	ls schnittstellen.LockSmith,
	of schnittstellen.ObjekteWriterFactory,
	af schnittstellen.AkteWriterFactory,
	reader TransactedReader[T3, *objekte.Transacted[T, T1, T2, T3, T4, T5]],
	delegate CreateOrUpdateDelegate[*objekte.Transacted[T, T1, T2, T3, T4, T5]],
	ma kennung.MatchableAdder,
	pmf persisted_metadatei_format.Format,
	kg konfig.Getter,
) (cou *createOrUpdate[T, T1, T2, T3, T4, T5]) {
	if pmf == nil {
		panic("nil persisted_metadatei_format.Format")
	}

	return &createOrUpdate[T, T1, T2, T3, T4, T5]{
		clock:                     clock,
		ls:                        ls,
		of:                        of,
		af:                        af,
		reader:                    reader,
		delegate:                  delegate,
		matchableAdder:            ma,
		persistentMetadateiFormat: pmf,
	}
}

func (cou createOrUpdate[T, T1, T2, T3, T4, T5]) CreateOrUpdateCheckedOut(
	co *objekte.CheckedOut[T, T1, T2, T3, T4, T5],
) (transactedPtr *objekte.Transacted[T, T1, T2, T3, T4, T5], err error) {
	kennungPtr := T3(&co.External.Sku.Kennung)
	objektePtr := T1(&co.External.Akte)

	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf("create or update %s", kennungPtr),
		}

		return
	}

	transactedPtr = &objekte.Transacted[T, T1, T2, T3, T4, T5]{
		Akte: *objektePtr,
		Sku: sku.Transacted[T2, T3]{
			Kennung: *kennungPtr,
			Schwanz: cou.clock.GetTai(),
		},
	}

	transactedPtr.SetAkteSha(co.External.GetAkteSha())

	var ow sha.WriteCloser

	if ow, err = cou.of.ObjekteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ow)

	if _, err = cou.persistentMetadateiFormat.FormatPersistentMetadatei(
		ow,
		transactedPtr,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	os := sha.Make(ow.Sha())
	transactedPtr.Sku.ObjekteSha = os

	if transactedPtr.GetObjekteSha().EqualsSha(co.Internal.GetObjekteSha()) {
		transactedPtr = &co.Internal

		if err = cou.delegate.Unchanged(transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = cou.matchableAdder.AddMatchable(transactedPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = cou.delegate.Updated(transactedPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (cou createOrUpdate[T, T1, T2, T3, T4, T5]) CreateOrUpdate(
	objektePtr T1,
	mg metadatei.Getter,
	kennungPtr T3,
) (transactedPtr *objekte.Transacted[T, T1, T2, T3, T4, T5], err error) {
	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf("create or update %s", kennungPtr.GetGattung()),
		}

		return
	}

	var mutter *objekte.Transacted[T, T1, T2, T3, T4, T5]

	if mutter, err = cou.reader.ReadOne(kennungPtr); err != nil {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	var m metadatei.Metadatei

	if mg != nil {
		m = mg.GetMetadatei()
	}

	transactedPtr = &objekte.Transacted[T, T1, T2, T3, T4, T5]{
		Metadatei: m,
		Akte:      *objektePtr,
		Sku: sku.Transacted[T2, T3]{
			Kennung: *kennungPtr,
			Schwanz: cou.clock.GetTai(),
		},
	}

	if mutter != nil {
		transactedPtr.Sku.Kopf = mutter.Sku.Kopf
	} else {
		errors.TodoP4("determine if this is necessary any more")
		// transactedPtr.Sku.Kopf = s.common.GetTransaktion().Time
	}

	objekte.CorrectAkteShaWith(transactedPtr, transactedPtr)

	var ow sha.WriteCloser

	if ow, err = cou.of.ObjekteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ow)

	if _, err = cou.persistentMetadateiFormat.FormatPersistentMetadatei(
		ow,
		transactedPtr,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	transactedPtr.Sku.ObjekteSha = sha.Make(ow.Sha())

	if mutter != nil &&
		transactedPtr.Sku.Kennung.Equals(mutter.Sku.Kennung) &&
		transactedPtr.GetObjekteSha().EqualsSha(mutter.GetObjekteSha()) {
		transactedPtr = mutter

		if err = cou.delegate.Unchanged(transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = cou.matchableAdder.AddMatchable(transactedPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter == nil {
		if err = cou.delegate.New(transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = cou.delegate.Updated(transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (cou createOrUpdate[T, T1, T2, T3, T4, T5]) CreateOrUpdateAkte(
	objektePtr T1,
	mg metadatei.Getter,
	kennungPtr T3,
	sh schnittstellen.Sha,
) (transactedPtr *objekte.Transacted[T, T1, T2, T3, T4, T5], err error) {
	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf("create or update %s", kennungPtr.GetGattung()),
		}

		return
	}

	var mutter *objekte.Transacted[T, T1, T2, T3, T4, T5]

	if mutter, err = cou.reader.ReadOne(kennungPtr); err != nil {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	var m metadatei.Metadatei

	if mg != nil {
		m = mg.GetMetadatei()
	}

	transactedPtr = &objekte.Transacted[T, T1, T2, T3, T4, T5]{
		Metadatei: m,
		Akte:      *objektePtr,
		Sku: sku.Transacted[T2, T3]{
			Kennung: *kennungPtr,
			Schwanz: cou.clock.GetTai(),
		},
	}

	transactedPtr.SetAkteSha(sh)

	if mutter != nil {
		transactedPtr.Sku.Kopf = mutter.Sku.Kopf
	} else {
		errors.TodoP4("determine if this is necessary any more")
		// transactedPtr.Sku.Kopf = s.common.GetTransaktion().Time
	}

	var ow sha.WriteCloser

	if ow, err = cou.of.ObjekteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ow)

	if _, err = cou.persistentMetadateiFormat.FormatPersistentMetadatei(
		ow,
		transactedPtr,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	transactedPtr.Sku.ObjekteSha = sha.Make(ow.Sha())

	if mutter != nil &&
		transactedPtr.Sku.Kennung.Equals(mutter.Sku.Kennung) &&
		transactedPtr.GetObjekteSha().EqualsSha(mutter.GetObjekteSha()) {
		transactedPtr = mutter

		if err = cou.delegate.Unchanged(transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = cou.matchableAdder.AddMatchable(transactedPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter == nil {
		if err = cou.delegate.New(transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = cou.delegate.Updated(transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
