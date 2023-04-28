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
	"github.com/friedenberg/zit/src/golf/persisted_metadatei_format"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type KonfigStore interface {
	reindexer

	GetAkteFormat() objekte.AkteFormat[erworben.Akte, *erworben.Akte]
	Read() (*erworben.Transacted, error)
	Update(*erworben.Akte) (*erworben.Transacted, error)

	objekte_store.TransactedLogger[*erworben.Transacted]
	objekte_store.AkteTextSaver[erworben.Akte, *erworben.Akte]
}

type KonfigInflator = objekte_store.TransactedInflator[
	erworben.Akte,
	*erworben.Akte,
	kennung.Konfig,
	*kennung.Konfig,
	objekte.NilVerzeichnisse[erworben.Akte],
	*objekte.NilVerzeichnisse[erworben.Akte],
]

type KonfigLogWriter = objekte_store.LogWriter[*erworben.Transacted]

type KonfigAkteTextSaver = objekte_store.AkteTextSaver[
	erworben.Akte,
	*erworben.Akte,
]

type konfigStore struct {
	store_util.StoreUtil

	schnittstellen.ObjekteIOFactory

	pool schnittstellen.Pool[erworben.Transacted, *erworben.Transacted]

	KonfigInflator
	KonfigAkteTextSaver
	KonfigLogWriter

	akteFormat objekte.AkteFormat[erworben.Akte, *erworben.Akte]
}

func (s *konfigStore) GetAkteFormat() objekte.AkteFormat[erworben.Akte, *erworben.Akte] {
	return s.akteFormat
}

func (s *konfigStore) SetLogWriter(
	tlw KonfigLogWriter,
) {
	s.KonfigLogWriter = tlw
}

func makeKonfigStore(
	sa store_util.StoreUtil,
) (s *konfigStore, err error) {
	pool := collections.MakePool[erworben.Transacted]()

	akteFormat := objekte_store.MakeAkteFormat[erworben.Akte, *erworben.Akte](
		objekte.MakeTextParserIgnoreTomlErrors[erworben.Akte](sa),
		objekte.ParsedAkteTomlFormatter[erworben.Akte]{},
		sa,
	)

	of := sa.ObjekteReaderWriterFactory(gattung.Konfig)

	s = &konfigStore{
		StoreUtil:        sa,
		ObjekteIOFactory: of,
		pool:             pool,
		KonfigInflator: objekte_store.MakeTransactedInflator[
			erworben.Akte,
			*erworben.Akte,
			kennung.Konfig,
			*kennung.Konfig,
			objekte.NilVerzeichnisse[erworben.Akte],
			*objekte.NilVerzeichnisse[erworben.Akte],
		](
			of,
			sa,
			persisted_metadatei_format.FormatForVersion(
				sa.GetKonfig().GetStoreVersion(),
			),
			akteFormat,
			pool,
		),
		KonfigAkteTextSaver: objekte_store.MakeAkteTextSaver[
			erworben.Akte,
			*erworben.Akte,
		](
			sa,
			akteFormat,
		),
		akteFormat: akteFormat,
	}

	return
}

func (s konfigStore) Flush() (err error) {
	return
}

func (s konfigStore) addOne(t *erworben.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().SetTransacted(t)
	return
}

func (s konfigStore) updateOne(t *erworben.Transacted) (err error) {
	s.StoreUtil.GetKonfigPtr().SetTransacted(t)
	return
}

func (s konfigStore) Update(
	ko *erworben.Akte,
) (kt *erworben.Transacted, err error) {
	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = errors.Wrap(objekte_store.ErrLockRequired{Operation: "update konfig"})
		return
	}

	var mutter *erworben.Transacted

	if mutter, err = s.Read(); err != nil {
		if errors.Is(err, objekte_store.ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	var akteSha schnittstellen.Sha

	if akteSha, _, err = s.SaveAkteText(*ko); err != nil {
		err = errors.Wrap(err)
		return
	}

	kt = &erworben.Transacted{
		Akte: *ko,
		Sku: sku.Transacted[kennung.Konfig, *kennung.Konfig]{
			Schwanz: s.StoreUtil.GetTai(),
		},
	}

	kt.SetAkteSha(akteSha)
	objekte.AssertAkteShasMatch(kt)

	// TODO-P3 refactor into reusable
	if mutter != nil {
		kt.Sku.Kopf = mutter.Sku.Kopf
	} else {
		kt.Sku.Kopf = s.StoreUtil.GetTai()
	}

	var ow sha.WriteCloser

	if ow, err = s.ObjekteIOFactory.ObjekteWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ow)

	if _, err = s.StoreUtil.GetPersistentMetadateiFormat().FormatPersistentMetadatei(
		ow,
		kt,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	kt.Sku.ObjekteSha = sha.Make(ow.Sha())
	mutterObjekteSha := mutter.GetObjekteSha()

	if mutter != nil && kt.GetObjekteSha().EqualsSha(mutterObjekteSha) {
		kt = mutter

		if err = s.KonfigLogWriter.Unchanged(kt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.StoreUtil.CommitTransacted(kt)
	s.StoreUtil.GetKonfigPtr().SetTransacted(kt)

	if err = s.StoreUtil.AddMatchable(kt); err != nil {
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
		Sku:  s.StoreUtil.GetKonfig().Sku,
		Akte: s.StoreUtil.GetKonfig().Akte,
	}

	if !tt.Sku.Schwanz.IsEmpty() {
		{
			var r sha.ReadCloser

			if r, err = s.ObjekteReader(
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

			if _, err = s.StoreUtil.GetPersistentMetadateiFormat().ParsePersistentMetadatei(
				r,
				tt,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		{
			var r sha.ReadCloser

			if r, err = s.ObjekteReader(
				tt.GetObjekteSha(),
			); err != nil {
				if errors.IsNotExist(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
				}

				return
			}

			defer errors.DeferredCloser(&err, r)

			fo := s.akteFormat

			var sh schnittstellen.Sha

			if sh, _, err = fo.ParseSaveAkte(r, &tt.Akte); err != nil {
				err = errors.Wrap(err)
				return
			}

			tt.SetAkteSha(sh)
		}
	}

	return
}

func (s *konfigStore) ReindexOne(
	sk sku.DataIdentity,
) (o kennung.Matchable, err error) {
	var te *erworben.Transacted
	defer s.pool.Put(te)

	if te, err = s.InflateFromDataIdentity(sk); err != nil {
		errors.Wrap(err)
		return
	}

	o = te

	s.KonfigLogWriter.Updated(te)

	return
}
