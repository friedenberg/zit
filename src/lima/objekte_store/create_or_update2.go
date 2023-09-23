package objekte_store

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/konfig"
)

type createOrUpdate2 struct {
	clock                     kennung.Clock
	ls                        schnittstellen.LockSmith
	af                        schnittstellen.AkteWriterFactory
	reader                    OneReader[*kennung.Kennung2, *sku.Transacted]
	delegate                  CreateOrUpdateDelegate[*sku.Transacted]
	matchableAdder            matcher.MatchableAdder
	persistentMetadateiFormat objekte_format.Format
	options                   objekte_format.Options
	kg                        konfig.Getter
	pool                      schnittstellen.Pool[sku.Transacted, *sku.Transacted]
}

func MakeCreateOrUpdate2(
	clock kennung.Clock,
	ls schnittstellen.LockSmith,
	af schnittstellen.AkteWriterFactory,
	reader OneReader[*kennung.Kennung2, *sku.Transacted],
	delegate CreateOrUpdateDelegate[*sku.Transacted],
	ma matcher.MatchableAdder,
	pmf objekte_format.Format,
	op objekte_format.Options,
	kg konfig.Getter,
	pool schnittstellen.Pool[sku.Transacted, *sku.Transacted],
) (cou *createOrUpdate2) {
	if pmf == nil {
		panic("nil persisted_metadatei_format.Format")
	}

	return &createOrUpdate2{
		clock:                     clock,
		ls:                        ls,
		af:                        af,
		reader:                    reader,
		delegate:                  delegate,
		matchableAdder:            ma,
		persistentMetadateiFormat: pmf,
		options:                   op,
		pool:                      pool,
	}
}

func (cou createOrUpdate2) CreateOrUpdateCheckedOut(
	co *objekte.CheckedOut2,
) (transactedPtr *sku.Transacted, err error) {
	kennungPtr := co.External.Kennung

	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf("create or update %s", kennungPtr),
		}

		return
	}

	transactedPtr = sku.GetTransactedPool().Get()
	transactedPtr.Kennung = kennungPtr
	transactedPtr.Metadatei = co.External.Metadatei
	transactedPtr.Metadatei.Tai = cou.clock.GetTai()

	transactedPtr.SetAkteSha(co.External.GetAkteSha())

	err = sku.CalculateAndSetSha(
		transactedPtr,
		cou.persistentMetadateiFormat,
		objekte_format.Options{IncludeTai: true},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO-P2: determine why Metadatei.Etiketten can be nil
	if transactedPtr.Metadatei.EqualsSansTai(co.Internal.Metadatei) &&
		co.State != checked_out_state.StateUntracked {
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

	if co.State == checked_out_state.StateUntracked {
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

func (cou createOrUpdate2) CreateOrUpdate(
	mg metadatei.Getter,
	kennungPtr *kennung.Kennung2,
) (transactedPtr *sku.Transacted, err error) {
	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var mutter *sku.Transacted

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

	transactedPtr = cou.pool.Get()

	transactedPtr.Kennung = *kennungPtr
	transactedPtr.Metadatei = m

	if mutter != nil {
		transactedPtr.Kopf = mutter.Kopf
	} else {
		errors.TodoP4("determine if this is necessary any more")
		// transactedPtr.Sku.Kopf = s.common.GetTransaktion().Time
	}

	if mutter != nil &&
		kennung.Equals(transactedPtr.GetKennung(), mutter.GetKennung()) &&
		transactedPtr.Metadatei.EqualsSansTai(mutter.Metadatei) {
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

func (cou createOrUpdate2) CreateOrUpdateAkte(
	mg metadatei.Getter,
	kennungPtr *kennung.Kennung2,
	sh schnittstellen.ShaLike,
) (transactedPtr *sku.Transacted, err error) {
	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf(
				"create or update %s",
				kennungPtr.GetGattung(),
			),
		}

		return
	}

	var mutter *sku.Transacted

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

	transactedPtr = sku.GetTransactedPool().Get()

	transactedPtr.Metadatei = m
	transactedPtr.Kennung = *kennungPtr
	transactedPtr.SetAkteSha(sh)

	err = sku.CalculateAndSetSha(
		transactedPtr,
		cou.persistentMetadateiFormat,
		objekte_format.Options{IncludeTai: true},
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter != nil {
		transactedPtr.Kopf = mutter.Kopf
	} else {
		errors.TodoP4("determine if this is necessary any more")
		// transactedPtr.Sku.Kopf = s.common.GetTransaktion().Time
	}

	if mutter != nil &&
		kennung.Equals(transactedPtr.GetKennung(), mutter.GetKennung()) &&
		transactedPtr.Metadatei.EqualsSansTai(mutter.Metadatei) {
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
