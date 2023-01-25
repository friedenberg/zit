package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/age_io"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/erworben"
)

type KonfigStore interface {
	reindexer
	GattungStore

	Read() (*erworben.Transacted, error)
	Update(*erworben.Objekte) (*erworben.Transacted, error)

	objekte.TransactedLogger[*erworben.Transacted]
	objekte.AkteTextSaver[erworben.Objekte, *erworben.Objekte]
}

type KonfigInflator = objekte.TransactedInflator[
	erworben.Objekte,
	*erworben.Objekte,
	kennung.Konfig,
	*kennung.Konfig,
	objekte.NilVerzeichnisse[erworben.Objekte],
	*objekte.NilVerzeichnisse[erworben.Objekte],
]

type KonfigLogWriter = objekte.LogWriter[*erworben.Transacted]

type KonfigAkteTextSaver = objekte.AkteTextSaver[
	erworben.Objekte,
	*erworben.Objekte,
]

type konfigStore struct {
	common *common

	pool collections.PoolLike[erworben.Transacted]

	KonfigInflator
	KonfigAkteTextSaver
	KonfigLogWriter
}

func (s *konfigStore) SetLogWriter(
	tlw KonfigLogWriter,
) {
	s.KonfigLogWriter = tlw
}

func makeKonfigStore(
	sa *common,
) (s *konfigStore, err error) {
	pool := collections.MakePool[erworben.Transacted]()

	s = &konfigStore{
		common: sa,
		pool:   pool,
		KonfigInflator: objekte.MakeTransactedInflator[
			erworben.Objekte,
			*erworben.Objekte,
			kennung.Konfig,
			*kennung.Konfig,
			objekte.NilVerzeichnisse[erworben.Objekte],
			*objekte.NilVerzeichnisse[erworben.Objekte],
		](
			sa,
			sa,
			nil,
			schnittstellen.Format[erworben.Objekte, *erworben.Objekte](
				erworben.MakeFormatText(sa),
			),
			pool,
		),
		KonfigAkteTextSaver: objekte.MakeAkteTextSaver[
			erworben.Objekte,
			*erworben.Objekte,
		](
			sa,
			&erworben.FormatterAkteTextToml{},
		),
	}

	return
}

func (s konfigStore) Flush() (err error) {
	return
}

func (s konfigStore) Update(
	ko *erworben.Objekte,
) (kt *erworben.Transacted, err error) {
	if !s.common.LockSmith.IsAcquired() {
		err = errors.Wrap(ErrLockRequired{Operation: "update konfig"})
		return
	}

	var w *age_io.Mover

	mo := age_io.MoveOptions{
		Age:                      s.common.Age,
		FinalPath:                s.common.GetStandort().DirObjektenKonfig(),
		GenerateFinalPathFromSha: true,
	}

	if w, err = age_io.NewMover(mo); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, w.Close)

	var mutter *erworben.Transacted

	if mutter, err = s.Read(); err != nil {
		if errors.Is(err, ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	kt = &erworben.Transacted{
		Objekte: *ko,
		Sku: sku.Transacted[kennung.Konfig, *kennung.Konfig]{
			Verzeichnisse: sku.Verzeichnisse{
				Schwanz: s.common.GetTransaktion().Time,
			},
		},
	}

	//TODO-P3 refactor into reusable
	if mutter != nil {
		kt.Sku.Kopf = mutter.Sku.Kopf
		kt.Sku.Mutter[0] = mutter.Sku.Schwanz
	} else {
		kt.Sku.Kopf = s.common.GetTransaktion().Time
	}

	fo := objekte.MakeFormat[erworben.Objekte, *erworben.Objekte]()

	if _, err = fo.Format(w, &kt.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	kt.Sku.ObjekteSha = sha.Make(w.Sha())

	if mutter != nil && kt.GetObjekteSha().Equals(mutter.GetObjekteSha()) {
		kt = mutter

		if err = s.KonfigLogWriter.Unchanged(kt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.common.GetTransaktion().Skus.Add(&kt.Sku)
	s.common.KonfigPtr().SetTransacted(kt)

	if err = s.common.Abbr.addStoredAbbreviation(kt); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.KonfigLogWriter.Updated(kt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s konfigStore) Read() (tt *erworben.Transacted, err error) {
	tt = &erworben.Transacted{
		Sku: s.common.Konfig().Sku,
		Objekte: erworben.Objekte{
			Akte: s.common.Konfig().Akte,
		},
	}

	if !tt.Sku.Schwanz.IsEmpty() {
		{
			var r sha.ReadCloser

			if r, err = s.common.ObjekteReader(
				gattung.Konfig,
				tt.Sku.ObjekteSha,
			); err != nil {
				if errors.IsNotExist(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			defer errors.Deferred(&err, r.Close)

			fo := objekte.MakeFormat[erworben.Objekte, *erworben.Objekte]()

			if _, err = fo.Parse(r, &tt.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		{
			var r sha.ReadCloser

			if r, err = s.common.ObjekteReader(
				gattung.Konfig,
				tt.Objekte.Sha,
			); err != nil {
				if errors.IsNotExist(err) {
					err = errors.Wrap(ErrNotFound{})
				} else {
					err = errors.Wrap(err)
				}
				return
			}

			defer errors.Deferred(&err, r.Close)

			fo := erworben.MakeFormatText(s.common)

			if _, err = fo.Parse(r, &tt.Objekte); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (s konfigStore) AllInChain() (c []*erworben.Transacted, err error) {
	return
}

func (s *konfigStore) reindexOne(
	sk sku.DataIdentity,
) (o schnittstellen.Stored, err error) {
	var te *erworben.Transacted
	defer s.pool.Put(te)

	if te, err = s.InflateFromDataIdentity(sk); err != nil {
		errors.Wrap(err)
		return
	}

	o = te

	s.common.KonfigPtr().SetTransacted(te)
	s.KonfigLogWriter.Updated(te)

	return
}
