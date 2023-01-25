package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
)

type KonfigStore interface {
	reindexer
	GattungStore

	Read() (*erworben.Transacted, error)
	Update(*erworben.Objekte) (*erworben.Transacted, error)

	objekte_store.TransactedLogger[*erworben.Transacted]
	objekte_store.AkteTextSaver[erworben.Objekte, *erworben.Objekte]
}

type KonfigInflator = objekte_store.TransactedInflator[
	erworben.Objekte,
	*erworben.Objekte,
	kennung.Konfig,
	*kennung.Konfig,
	objekte.NilVerzeichnisse[erworben.Objekte],
	*objekte.NilVerzeichnisse[erworben.Objekte],
]

type KonfigLogWriter = objekte_store.LogWriter[*erworben.Transacted]

type KonfigAkteTextSaver = objekte_store.AkteTextSaver[
	erworben.Objekte,
	*erworben.Objekte,
]

type konfigStore struct {
	StoreUtil

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
	sa StoreUtil,
) (s *konfigStore, err error) {
	pool := collections.MakePool[erworben.Transacted]()

	s = &konfigStore{
		StoreUtil: sa,
		pool:      pool,
		KonfigInflator: objekte_store.MakeTransactedInflator[
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
		KonfigAkteTextSaver: objekte_store.MakeAkteTextSaver[
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
	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = errors.Wrap(objekte_store.ErrLockRequired{Operation: "update konfig"})
		return
	}

	var ow sha.WriteCloser

	if ow, err = s.StoreUtil.ObjekteWriter(gattung.Konfig); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ow)

	var mutter *erworben.Transacted

	if mutter, err = s.Read(); err != nil {
		if errors.Is(err, objekte_store.ErrNotFound{}) {
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
				Schwanz: s.StoreUtil.GetTime(),
			},
		},
	}

	//TODO-P3 refactor into reusable
	if mutter != nil {
		kt.Sku.Kopf = mutter.Sku.Kopf
		kt.Sku.Mutter[0] = mutter.Sku.Schwanz
	} else {
		kt.Sku.Kopf = s.StoreUtil.GetTime()
	}

	fo := objekte.MakeFormat[erworben.Objekte, *erworben.Objekte]()

	if _, err = fo.Format(ow, &kt.Objekte); err != nil {
		err = errors.Wrap(err)
		return
	}

	kt.Sku.ObjekteSha = sha.Make(ow.Sha())

	if mutter != nil && kt.GetObjekteSha().Equals(mutter.GetObjekteSha()) {
		kt = mutter

		if err = s.KonfigLogWriter.Unchanged(kt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.StoreUtil.CommitTransacted(kt)
	s.StoreUtil.GetKonfigPtr().SetTransacted(kt)

	if err = s.StoreUtil.GetAbbrStore().addStoredAbbreviation(kt); err != nil {
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
		Sku: s.StoreUtil.GetKonfig().Sku,
		Objekte: erworben.Objekte{
			Akte: s.StoreUtil.GetKonfig().Akte,
		},
	}

	if !tt.Sku.Schwanz.IsEmpty() {
		{
			var r sha.ReadCloser

			if r, err = s.StoreUtil.ObjekteReader(
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

			if r, err = s.StoreUtil.ObjekteReader(
				gattung.Konfig,
				tt.Objekte.Sha,
			); err != nil {
				if errors.IsNotExist(err) {
					err = errors.Wrap(objekte_store.ErrNotFound{})
				} else {
					err = errors.Wrap(err)
				}
				return
			}

			defer errors.Deferred(&err, r.Close)

			fo := erworben.MakeFormatText(s.StoreUtil)

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

	s.StoreUtil.GetKonfigPtr().SetTransacted(te)
	s.KonfigLogWriter.Updated(te)

	return
}
