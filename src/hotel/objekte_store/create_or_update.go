package objekte_store

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/echo/ts"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
)

type CreateOrUpdateDelegate[T any] struct {
	New       collections.WriterFunc[T]
	Updated   collections.WriterFunc[T]
	Unchanged collections.WriterFunc[T]
}

type createOrUpdate[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
	T4 any,
	T5 schnittstellen.VerzeichnissePtr[T4, T],
] struct {
	clock    ts.Clock
	ls       schnittstellen.LockSmith
	oaf      schnittstellen.ObjekteAkteWriterFactory
	reader   TransactedReader[T3, *objekte.Transacted[T, T1, T2, T3, T4, T5]]
	delegate CreateOrUpdateDelegate[*objekte.Transacted[T, T1, T2, T3, T4, T5]]
}

func MakeCreateOrUpdate[
	T schnittstellen.Objekte[T],
	T1 schnittstellen.ObjektePtr[T],
	T2 schnittstellen.Id[T2],
	T3 schnittstellen.IdPtr[T2],
	T4 any,
	T5 schnittstellen.VerzeichnissePtr[T4, T],
](
	clock ts.Clock,
	ls schnittstellen.LockSmith,
	oaf schnittstellen.ObjekteAkteWriterFactory,
	reader TransactedReader[T3, *objekte.Transacted[T, T1, T2, T3, T4, T5]],
	delegate CreateOrUpdateDelegate[*objekte.Transacted[T, T1, T2, T3, T4, T5]],
) (cou *createOrUpdate[T, T1, T2, T3, T4, T5]) {
	return &createOrUpdate[T, T1, T2, T3, T4, T5]{
		clock:    clock,
		ls:       ls,
		oaf:      oaf,
		reader:   reader,
		delegate: delegate,
	}
}

func (cou createOrUpdate[T, T1, T2, T3, T4, T5]) CreateOrUpdate(
	objektePtr T1,
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

	transactedPtr = &objekte.Transacted[T, T1, T2, T3, T4, T5]{
		Objekte: *objektePtr,
		Sku: sku.Transacted[T2, T3]{
			Kennung: *kennungPtr,
			Verzeichnisse: sku.Verzeichnisse{
				Schwanz: cou.clock.GetTime(),
			},
		},
	}

	if mutter != nil {
		transactedPtr.Sku.Kopf = mutter.Sku.Kopf
		transactedPtr.Sku.Mutter[0] = mutter.Sku.Schwanz
	} else {
		errors.TodoP4("determine if this is necessary any more")
		// transactedPtr.Sku.Kopf = s.common.GetTransaktion().Time
	}

	fo := objekte.MakeFormat[T, T1]()

	var ow sha.WriteCloser

	if ow, err = cou.oaf.ObjekteWriter(kennungPtr.GetGattung()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ow)

	if _, err = fo.Format(ow, &transactedPtr.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	transactedPtr.Sku.ObjekteSha = sha.Make(ow.Sha())

	if mutter != nil && transactedPtr.GetObjekteSha().Equals(mutter.GetObjekteSha()) {
		transactedPtr = mutter

		if err = cou.delegate.Unchanged(transactedPtr); err != nil {
			err = errors.Wrap(err)
			return
		}

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