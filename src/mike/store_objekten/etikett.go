package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/age_io"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/india/etikett"
	"github.com/friedenberg/zit/src/schnittstellen"
)

type EtikettStore interface {
	reindexer
	GattungStore

	objekte.TransactedLogger[*etikett.Transacted]

	objekte.TransactedReader[
		*kennung.Etikett,
		*etikett.Transacted,
	]

	objekte.CreateOrUpdater[
		*etikett.Objekte,
		*kennung.Etikett,
		*etikett.Transacted,
	]
}

type EtikettInflator = objekte.TransactedInflator[
	etikett.Objekte,
	*etikett.Objekte,
	kennung.Etikett,
	*kennung.Etikett,
	objekte.NilVerzeichnisse[etikett.Objekte],
	*objekte.NilVerzeichnisse[etikett.Objekte],
]

type EtikettLogWriter = objekte.LogWriter[*etikett.Transacted]

type EtikettAkteTextSaver = objekte.AkteTextSaver[
	etikett.Objekte,
	*etikett.Objekte,
]

type etikettStore struct {
	common *common

	pool collections.PoolLike[etikett.Transacted]

	EtikettInflator
	EtikettAkteTextSaver
	EtikettLogWriter
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
		EtikettInflator: objekte.MakeTransactedInflator[
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
		EtikettAkteTextSaver: objekte.MakeAkteTextSaver[
			etikett.Objekte,
			*etikett.Objekte,
		](
			sa,
			&etikett.FormatterAkteTextToml{},
		),
	}

	return
}

func (s etikettStore) Flush() (err error) {
	return
}

func (s etikettStore) CreateOrUpdate(
	to *etikett.Objekte,
	tk *kennung.Etikett,
) (tt *etikett.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = ErrLockRequired{
			Operation: "create or update etikett",
		}

		return
	}

	var mutter *etikett.Transacted

	if mutter, err = s.ReadOne(tk); err != nil {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	tt = &etikett.Transacted{
		Objekte: *to,
		Sku: sku.Transacted[kennung.Etikett, *kennung.Etikett]{
			Kennung: *tk,
			Verzeichnisse: sku.Verzeichnisse{
				Schwanz: s.common.Transaktion.Time,
			},
		},
	}

	//TODO-P3 refactor into reusable
	if mutter != nil {
		tt.Sku.Kopf = mutter.Sku.Kopf
		tt.Sku.Mutter[0] = mutter.Sku.Schwanz
	} else {
		tt.Sku.Kopf = s.common.Transaktion.Time
	}

	fo := objekte.MakeFormat[etikett.Objekte, *etikett.Objekte]()

	var w *age_io.Mover

	mo := age_io.MoveOptions{
		Age:                      s.common.Age,
		FinalPath:                s.common.Standort.DirObjektenEtiketten(),
		GenerateFinalPathFromSha: true,
	}

	if w, err = age_io.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w.Close)

	if _, err = fo.Format(w, &tt.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	tt.Sku.ObjekteSha = sha.Make(w.Sha())

	if mutter != nil && tt.GetObjekteSha().Equals(mutter.GetObjekteSha()) {
		tt = mutter

		if err = s.EtikettLogWriter.Unchanged(tt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.common.Transaktion.Skus.Add(&tt.Sku)
	s.common.KonfigPtr().AddEtikett(tt)

	if mutter == nil {
		if err = s.EtikettLogWriter.New(tt); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = s.EtikettLogWriter.Updated(tt); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s etikettStore) ReadOne(
	k *kennung.Etikett,
) (tt *etikett.Transacted, err error) {
	tt = s.common.Konfig().GetEtikett(*k)

	if tt == nil {
		err = errors.Wrap(ErrNotFound{Id: k})
		return
	}

	return
}

func (s etikettStore) ReadAllSchwanzen(
	f collections.WriterFunc[*etikett.Transacted],
) (err error) {
	//TODO-P2
	return
}

func (s etikettStore) ReadAll(
	f collections.WriterFunc[*etikett.Transacted],
) (err error) {
	//TODO-P2
	return
}

func (s etikettStore) AllInChain(k kennung.Etikett) (c []*etikett.Transacted, err error) {
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
