package objekte_store

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/konfig"
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
	reader                    TransactedReader[T3, sku.SkuLikePtr]
	delegate                  CreateOrUpdateDelegate[*sku.Transacted[T2, T3]]
	matchableAdder            matcher.MatchableAdder
	persistentMetadateiFormat objekte_format.Format
	options                   objekte_format.Options
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
	reader TransactedReader[T3, sku.SkuLikePtr],
	delegate CreateOrUpdateDelegate[*sku.Transacted[T2, T3]],
	ma matcher.MatchableAdder,
	pmf objekte_format.Format,
	op objekte_format.Options,
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
		options:                   op,
	}
}

func (cou createOrUpdate[T, T1, T2, T3]) CreateOrUpdateCheckedOut(
	co *objekte.CheckedOut[T2, T3],
) (transactedPtr *sku.Transacted[T2, T3], err error) {
	kennungPtr := T3(&co.External.Kennung)

	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf("create or update %s", kennungPtr),
		}

		return
	}

	transactedPtr = &sku.Transacted[T2, T3]{
		Kennung: *kennungPtr,
		Metadatei: metadatei.Metadatei{
			Tai: cou.clock.GetTai(),
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
		cou.options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	os := sha.Make(ow.GetShaLike())
	transactedPtr.ObjekteSha = os

	// TODO-P2: determine why Metadatei.Etiketten can be nil
	if transactedPtr.Metadatei.EqualsSansTai(co.Internal.Metadatei) {
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
	mg metadatei.Getter,
	kennungPtr T3,
) (transactedPtr *sku.Transacted[T2, T3], err error) {
	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var mutter *sku.Transacted2

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

	m.Tai = cou.clock.GetTai()

	transactedPtr = &sku.Transacted[T2, T3]{
		Kennung:   *kennungPtr,
		Metadatei: m,
	}

	if mutter != nil {
		transactedPtr.Kopf = mutter.GetKopf()
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
		cou.options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	transactedPtr.ObjekteSha = sha.Make(ow.GetShaLike())

	if mutter != nil &&
		kennung.Equals(transactedPtr.GetKennung(), mutter.GetKennungLike()) &&
		transactedPtr.Metadatei.EqualsSansTai(mutter.GetMetadatei()) {
		if err = transactedPtr.SetFromSkuLike(mutter); err != nil {
			err = errors.Wrap(err)
			return
		}

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
	mg metadatei.Getter,
	kennungPtr T3,
	sh schnittstellen.ShaLike,
) (transactedPtr *sku.Transacted[T2, T3], err error) {
	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var mutter sku.SkuLikePtr

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

	m.Tai = cou.clock.GetTai()

	transactedPtr = &sku.Transacted[T2, T3]{
		Metadatei: m,
		Kennung:   *kennungPtr,
	}

	transactedPtr.SetAkteSha(sh)

	if mutter != nil {
		transactedPtr.Kopf = mutter.GetKopf()
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
		cou.options,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	transactedPtr.ObjekteSha = sha.Make(ow.GetShaLike())

	if mutter != nil &&
		kennung.Equals(transactedPtr.GetKennung(), mutter.GetKennungLike()) &&
		transactedPtr.Metadatei.EqualsSansTai(mutter.GetMetadatei()) {
		if err = transactedPtr.SetFromSkuLike(mutter); err != nil {
			err = errors.Wrap(err)
			return
		}

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
