package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/etikett"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
)

type EtikettStore interface {
	reindexer
	GattungStore

	objekte_store.TransactedLogger[*etikett.Transacted]

	objekte_store.TransactedReader[
		*kennung.Etikett,
		*etikett.Transacted,
	]

	objekte_store.CreateOrUpdater[
		*etikett.Objekte,
		*kennung.Etikett,
		*etikett.Transacted,
	]
}

type EtikettTransactedReader = objekte_store.TransactedReader[
	*kennung.Etikett,
	*etikett.Transacted,
]

type EtikettInflator = objekte_store.TransactedInflator[
	etikett.Objekte,
	*etikett.Objekte,
	kennung.Etikett,
	*kennung.Etikett,
	objekte.NilVerzeichnisse[etikett.Objekte],
	*objekte.NilVerzeichnisse[etikett.Objekte],
]

type EtikettLogWriter = objekte_store.LogWriter[*etikett.Transacted]

type EtikettAkteTextSaver = objekte_store.AkteTextSaver[
	etikett.Objekte,
	*etikett.Objekte,
]

type etikettStore struct {
	common *common

	pool collections.PoolLike[etikett.Transacted]

	EtikettInflator
	EtikettAkteTextSaver
	EtikettLogWriter

	objekte_store.CreateOrUpdater[
		*etikett.Objekte,
		*kennung.Etikett,
		*etikett.Transacted,
	]
}

func (s *etikettStore) SetLogWriter(
	tlw EtikettLogWriter,
) {
	s.EtikettLogWriter = tlw
}

func makeEtikettStore(
	sa *common,
) (s *etikettStore, err error) {
	pool := collections.MakePool[etikett.Transacted]()

	s = &etikettStore{
		common: sa,
		pool:   pool,
		EtikettInflator: objekte_store.MakeTransactedInflator[
			etikett.Objekte,
			*etikett.Objekte,
			kennung.Etikett,
			*kennung.Etikett,
			objekte.NilVerzeichnisse[etikett.Objekte],
			*objekte.NilVerzeichnisse[etikett.Objekte],
		](
			sa,
			sa,
			nil,
			schnittstellen.Format[etikett.Objekte, *etikett.Objekte](
				etikett.MakeFormatText(sa),
			),
			pool,
		),
		EtikettAkteTextSaver: objekte_store.MakeAkteTextSaver[
			etikett.Objekte,
			*etikett.Objekte,
		](
			sa,
			&etikett.FormatterAkteTextToml{},
		),
	}

	newOrUpdated := func(t *etikett.Transacted) (err error) {
		s.common.AddSku(t.GetSku2(), &t.Sku)
		s.common.KonfigPtr().AddEtikett(t)

		return
	}

	s.CreateOrUpdater = objekte_store.MakeCreateOrUpdate(
		sa,
		sa.LockSmith,
		sa,
		EtikettTransactedReader(s),
		objekte_store.CreateOrUpdateDelegate[*etikett.Transacted]{
			New: func(t *etikett.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.EtikettLogWriter.New(t)
			},
			Updated: func(t *etikett.Transacted) (err error) {
				if err = newOrUpdated(t); err != nil {
					err = errors.Wrap(err)
					return
				}

				return s.EtikettLogWriter.Updated(t)
			},
			Unchanged: func(t *etikett.Transacted) (err error) {
				return s.EtikettLogWriter.Unchanged(t)
			},
		},
	)

	return
}

func (s etikettStore) Flush() (err error) {
	return
}

func (s etikettStore) ReadOne(
	k *kennung.Etikett,
) (tt *etikett.Transacted, err error) {
	tt = s.common.Konfig().GetEtikett(*k)

	if tt == nil {
		err = errors.Wrap(objekte_store.ErrNotFound{Id: k})
		return
	}

	return
}

func (s etikettStore) ReadAllSchwanzen(
	f collections.WriterFunc[*etikett.Transacted],
) (err error) {
	errors.TodoP2("implement")
	return
}

func (s etikettStore) ReadAll(
	f collections.WriterFunc[*etikett.Transacted],
) (err error) {
	errors.TodoP2("implement")
	return
}

func (s *etikettStore) reindexOne(
	sk sku.DataIdentity,
) (o schnittstellen.Stored, err error) {
	var te *etikett.Transacted
	defer s.pool.Put(te)

	if te, err = s.InflateFromDataIdentity(sk); err != nil {
		errors.Wrap(err)
		return
	}

	o = te

	s.common.KonfigPtr().AddEtikett(te)

	if te.IsNew() {
		s.EtikettLogWriter.New(te)
	} else {
		s.EtikettLogWriter.Updated(te)
	}

	return
}
