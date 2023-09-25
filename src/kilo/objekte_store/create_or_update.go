package objekte_store

import (
	"fmt"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/juliett/konfig"
)

type CreateOrUpdateDelegate struct {
	New       schnittstellen.FuncIter[*sku.Transacted]
	Updated   schnittstellen.FuncIter[*sku.Transacted]
	Unchanged schnittstellen.FuncIter[*sku.Transacted]
}

type createOrUpdate struct {
	clock                     kennung.Clock
	ls                        schnittstellen.LockSmith
	af                        schnittstellen.AkteWriterFactory
	reader                    OneReader
	delegate                  CreateOrUpdateDelegate
	matchableAdder            matcher.MatchableAdder
	persistentMetadateiFormat objekte_format.Format
	options                   objekte_format.Options
	kg                        konfig.Getter
}

func MakeCreateOrUpdate(
	clock kennung.Clock,
	ls schnittstellen.LockSmith,
	af schnittstellen.AkteWriterFactory,
	reader OneReader,
	delegate CreateOrUpdateDelegate,
	ma matcher.MatchableAdder,
	pmf objekte_format.Format,
	op objekte_format.Options,
	kg konfig.Getter,
) (cou *createOrUpdate) {
	if pmf == nil {
		panic("nil persisted_metadatei_format.Format")
	}

	return &createOrUpdate{
		clock:                     clock,
		ls:                        ls,
		af:                        af,
		reader:                    reader,
		delegate:                  delegate,
		matchableAdder:            ma,
		persistentMetadateiFormat: pmf,
		options:                   op,
	}
}

func (cou createOrUpdate) CreateOrUpdateCheckedOut(
	co *sku.CheckedOut,
) (transactedPtr *sku.Transacted, err error) {
	kennungPtr := co.External.Kennung

	if !cou.ls.IsAcquired() {
		err = ErrLockRequired{
			Operation: fmt.Sprintf("create or update %s", kennungPtr),
		}

		return
	}

	transactedPtr = sku.GetTransactedPool().Get()

	if err = transactedPtr.SetFromSkuLike(&co.External); err != nil {
		err = errors.Wrap(err)
		return
	}

	transactedPtr.Metadatei.Tai = cou.clock.GetTai()
	transactedPtr.SetAkteSha(co.External.GetAkteSha())

	err = sku.CalculateAndSetSha(
		transactedPtr,
		cou.persistentMetadateiFormat,
		cou.options,
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

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

func (cou createOrUpdate) CreateOrUpdate(
	mg metadatei.Getter,
	kennungPtr kennung.Kennung,
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

	if err = transactedPtr.Kennung.SetWithKennung(kennungPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	if mutter != nil {
		transactedPtr.Kopf = mutter.GetKopf()
	} else {
		errors.TodoP4("determine if this is necessary any more")
		// transactedPtr.Sku.Kopf = s.common.GetTransaktion().Time
	}

	err = sku.CalculateAndSetSha(
		transactedPtr,
		cou.persistentMetadateiFormat,
		cou.options,
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

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

func (cou createOrUpdate) CreateOrUpdateAkte(
	mg metadatei.Getter,
	kennungPtr kennung.Kennung,
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

	if err = transactedPtr.Kennung.SetWithKennung(kennungPtr); err != nil {
		err = errors.Wrap(err)
		return
	}

	transactedPtr.SetAkteSha(sh)

	if mutter != nil {
		transactedPtr.Kopf = mutter.GetKopf()
	} else {
		errors.TodoP4("determine if this is necessary any more")
		// transactedPtr.Sku.Kopf = s.common.GetTransaktion().Time
	}

	err = sku.CalculateAndSetSha(
		transactedPtr,
		cou.persistentMetadateiFormat,
		cou.options,
	)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

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
