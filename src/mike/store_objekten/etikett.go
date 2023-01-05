package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/golf/age_io"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/objekte"
	"github.com/friedenberg/zit/src/india/etikett"
)

type etikettLogWriter = collections.WriterFunc[*etikett.Transacted]

type EtikettLogWriters struct {
	New, Updated, Unchanged etikettLogWriter
}

type etikettStore struct {
	common *common

	pool collections.PoolLike[etikett.Transacted]

	objekte.TransactedInflator[
		etikett.Objekte,
		*etikett.Objekte,
		kennung.Etikett,
		*kennung.Etikett,
		objekte.NilVerzeichnisse[etikett.Objekte],
		*objekte.NilVerzeichnisse[etikett.Objekte],
	]

	objekte.AkteTextSaver[
		etikett.Objekte,
		*etikett.Objekte,
	]

	EtikettLogWriters
}

func (s *etikettStore) SetEtikettLogWriters(
	tlw EtikettLogWriters,
) {
	s.EtikettLogWriters = tlw
}

func makeEtikettStore(
	sa *common,
) (s *etikettStore, err error) {
	pool := collections.MakePool[etikett.Transacted]()

	s = &etikettStore{
		common: sa,
		pool:   pool,
		TransactedInflator: objekte.MakeTransactedInflator[
			etikett.Objekte,
			*etikett.Objekte,
			kennung.Etikett,
			*kennung.Etikett,
			objekte.NilVerzeichnisse[etikett.Objekte],
			*objekte.NilVerzeichnisse[etikett.Objekte],
		](
			sa.ReadCloserObjektenSku,
			sa.AkteReader,
			nil,
			gattung.Parser[etikett.Objekte, *etikett.Objekte](
				etikett.MakeFormatText(sa),
			),
			pool,
		),
		AkteTextSaver: objekte.MakeAkteTextSaver[
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
			Schwanz: s.common.Transaktion.Time,
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

	tt.Sku.Sha = w.Sha()

	if mutter != nil && tt.ObjekteSha().Equals(mutter.ObjekteSha()) {
		tt = mutter

		if err = s.EtikettLogWriters.Unchanged(tt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.common.Transaktion.Skus.Add(&tt.Sku)
	s.common.KonfigPtr().AddEtikett(tt)

	if mutter == nil {
		if err = s.EtikettLogWriters.New(tt); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = s.EtikettLogWriters.Updated(tt); err != nil {
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

func (s etikettStore) AllInChain(k kennung.Etikett) (c []*etikett.Transacted, err error) {
	return
}

func (s *etikettStore) reindexOne(
	t *transaktion.Transaktion,
	o sku.SkuLike,
) (err error) {
	var te *etikett.Transacted
	defer s.pool.Put(te)

	if te, err = s.Inflate(t.Time, o); err != nil {
		errors.Wrap(err)
		return
	}

	s.common.KonfigPtr().AddEtikett(te)

	if te.IsNew() {
		s.EtikettLogWriters.New(te)
	} else {
		s.EtikettLogWriters.Updated(te)
	}

	return
}
