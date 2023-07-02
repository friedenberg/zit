package objekte_store

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/objekte_format"
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
] struct {
	clock                     kennung.Clock
	ls                        schnittstellen.LockSmith
	of                        schnittstellen.ObjekteWriterFactory
	af                        schnittstellen.AkteWriterFactory
	reader                    TransactedReader[T3, *objekte.Transacted[T, T1, T2, T3]]
	delegate                  CreateOrUpdateDelegate[*objekte.Transacted[T, T1, T2, T3]]
	matchableAdder            kennung.MatchableAdder
	persistentMetadateiFormat objekte_format.Format
	kg                        konfig.Getter
}

func MakeCreateOrUpdate[
	T objekte.Akte[T],
	T1 objekte.AktePtr[T],
	T2 kennung.KennungLike[T2],
	T3 kennung.KennungLikePtr[T2],
](
	clock kennung.Clock,
	ls schnittstellen.LockSmith,
	of schnittstellen.ObjekteWriterFactory,
	af schnittstellen.AkteWriterFactory,
	reader TransactedReader[T3, *objekte.Transacted[T, T1, T2, T3]],
	delegate CreateOrUpdateDelegate[*objekte.Transacted[T, T1, T2, T3]],
	ma kennung.MatchableAdder,
	pmf objekte_format.Format,
	kg konfig.Getter,
) (cou *createOrUpdate[T, T1, T2, T3]) {
	if pmf == nil {
		panic("nil persisted_metadatei_format.Format")
	}

	return &createOrUpdate[T, T1, T2, T3]{
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

func (cou createOrUpdate[T, T1, T2, T3]) CreateOrUpdateCheckedOut(
	co *objekte.CheckedOut[T, T1, T2, T3],
) (transactedPtr *objekte.Transacted[T, T1, T2, T3], err error) {
	kennungPtr := T3(&co.External.Sku.WithKennung.Kennung)
	objektePtr := T1(&co.External.Akte)

	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf("create or update %s", kennungPtr),
		}

		return
	}

	transactedPtr = &objekte.Transacted[T, T1, T2, T3]{
		Akte: *objektePtr,
		Sku: sku.Transacted[T2, T3]{
			WithKennung: sku.WithKennung[T2, T3]{
				Kennung: *kennungPtr,
				Metadatei: sku.Metadatei{
					Tai: cou.clock.GetTai(),
				},
			},
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

	os := sha.Make(ow.GetShaLike())
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

func (cou createOrUpdate[T, T1, T2, T3]) CreateOrUpdate(
	objektePtr T1,
	mg metadatei.Getter,
	kennungPtr T3,
) (transactedPtr *objekte.Transacted[T, T1, T2, T3], err error) {
	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var mutter *objekte.Transacted[T, T1, T2, T3]

	if mutter, err = cou.reader.ReadOne(kennungPtr); err != nil {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	var m sku.Metadatei

	if mg != nil {
		m = mg.GetMetadatei()
	}

	m.Tai = cou.clock.GetTai()

	transactedPtr = &objekte.Transacted[T, T1, T2, T3]{
		Akte: *objektePtr,
		Sku: sku.Transacted[T2, T3]{
			WithKennung: sku.WithKennung[T2, T3]{
				Kennung:   *kennungPtr,
				Metadatei: m,
			},
		},
	}

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

	transactedPtr.Sku.ObjekteSha = sha.Make(ow.GetShaLike())

	if mutter != nil &&
		transactedPtr.Sku.GetKennung().Equals(mutter.Sku.GetKennung()) &&
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

func (cou createOrUpdate[T, T1, T2, T3]) CreateOrUpdateAkte(
	objektePtr T1,
	mg metadatei.Getter,
	kennungPtr T3,
	sh schnittstellen.ShaLike,
) (transactedPtr *objekte.Transacted[T, T1, T2, T3], err error) {
	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var mutter *objekte.Transacted[T, T1, T2, T3]

	if mutter, err = cou.reader.ReadOne(kennungPtr); err != nil {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	var m sku.Metadatei

	if mg != nil {
		m = mg.GetMetadatei()
	}

	m.Tai = cou.clock.GetTai()

	transactedPtr = &objekte.Transacted[T, T1, T2, T3]{
		Akte: *objektePtr,
		Sku: sku.Transacted[T2, T3]{
			WithKennung: sku.WithKennung[T2, T3]{
				Metadatei: m,
				Kennung:   *kennungPtr,
			},
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

	transactedPtr.Sku.ObjekteSha = sha.Make(ow.GetShaLike())

	if mutter != nil &&
		transactedPtr.Sku.GetKennung().Equals(mutter.Sku.GetKennung()) &&
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
