package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/foxtrot/sku"
	"github.com/friedenberg/zit/src/golf/objekte"
	"github.com/friedenberg/zit/src/golf/transaktion"
	"github.com/friedenberg/zit/src/hotel/erworben"
	"github.com/friedenberg/zit/src/hotel/objekte_store"
	"github.com/friedenberg/zit/src/india/bestandsaufnahme"
	"github.com/friedenberg/zit/src/kilo/store_util"
)

type KonfigStore interface {
	GetAkteFormat() objekte.AkteFormat[erworben.Akte, *erworben.Akte]
	Update(*erworben.Akte, schnittstellen.ShaLike) (*erworben.Transacted, error)

	CommonStoreBase[
		erworben.Akte,
		*erworben.Akte,
		kennung.Konfig,
		*kennung.Konfig,
	]
}

type konfigStore struct {
	*commonStore[
		erworben.Akte,
		*erworben.Akte,
		kennung.Konfig,
		*kennung.Konfig,
	]

	akteFormat objekte.AkteFormat[erworben.Akte, *erworben.Akte]
}

func (s *konfigStore) GetAkteFormat() objekte.AkteFormat[erworben.Akte, *erworben.Akte] {
	return s.akteFormat
}

func makeKonfigStore(
	sa store_util.StoreUtil,
) (s *konfigStore, err error) {
	s = &konfigStore{
		akteFormat: objekte_store.MakeAkteFormat[erworben.Akte, *erworben.Akte](
			objekte.MakeTextParserIgnoreTomlErrors[erworben.Akte](sa),
			objekte.ParsedAkteTomlFormatter[erworben.Akte]{},
			sa,
		),
	}

	s.commonStore, err = makeCommonStore[
		erworben.Akte,
		*erworben.Akte,
		kennung.Konfig,
		*kennung.Konfig,
	](
		gattung.Konfig,
		s,
		sa,
		s,
		s.akteFormat,
	)

	if s.commonStore.ObjekteSaver == nil {
		panic("ObjekteSaver is nil")
	}

	if err != nil {
		err = errors.Wrap(err)
		return
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
	sh schnittstellen.ShaLike,
) (kt *erworben.Transacted, err error) {
	if !s.StoreUtil.GetLockSmith().IsAcquired() {
		err = errors.Wrap(
			objekte_store.ErrLockRequired{Operation: "update konfig"},
		)
		return
	}

	var mutter *erworben.Transacted

	if mutter, err = s.ReadOne(&kennung.Konfig{}); err != nil {
		if errors.Is(err, objekte_store.ErrNotFound{}) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	kt = &erworben.Transacted{
		Akte: *ko,
	}

	kt.SetTai(s.StoreUtil.GetTai())
	kt.SetAkteSha(sh)

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

	kt.Sku.ObjekteSha = sha.Make(ow.GetShaLike())
	mutterObjekteSha := mutter.GetObjekteSha()

	if mutter != nil && kt.GetObjekteSha().EqualsSha(mutterObjekteSha) {
		kt = mutter

		if err = s.LogWriter.Unchanged(kt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	s.StoreUtil.CommitUpdatedTransacted(kt)
	s.StoreUtil.GetKonfigPtr().SetTransacted(kt)

	if err = s.StoreUtil.AddMatchable(kt); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.LogWriter.Updated(kt); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *konfigStore) ReadAllSchwanzen(
	w schnittstellen.FuncIter[*erworben.Transacted],
) (err error) {
	var k *erworben.Transacted

	if k, err = i.ReadOne(&kennung.Konfig{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = w(k); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *konfigStore) ReadAll(
	w schnittstellen.FuncIter[*erworben.Transacted],
) (err error) {
	if s.StoreUtil.GetKonfig().UseBestandsaufnahme {
		f1 := func(t *bestandsaufnahme.Transacted) (err error) {
			if err = sku.HeapEach(
				t.Akte.Skus,
				func(sk sku.SkuLike) (err error) {
					if sk.GetGattung() != gattung.Konfig {
						return
					}

					var te *erworben.Transacted

					if te, err = s.InflateFromDataIdentity(sk); err != nil {
						if errors.Is(err, toml.Error{}) {
							err = nil
						} else {
							err = errors.Wrap(err)
							return
						}
					}

					if err = w(te); err != nil {
						err = errors.Wrap(err)
						return
					}

					return
				},
			); err != nil {
				err = errors.Wrapf(
					err,
					"Bestandsaufnahme: %s",
					t.GetKennung(),
				)

				return
			}

			return
		}

		if err = s.StoreUtil.GetBestandsaufnahmeStore().ReadAll(f1); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = s.StoreUtil.GetTransaktionStore().ReadAllTransaktions(
			func(t *transaktion.Transaktion) (err error) {
				if err = t.Skus.Each(
					func(o sku.SkuLike) (err error) {
						if o.GetGattung() != gattung.Konfig {
							return
						}

						var te *erworben.Transacted

						if te, err = s.InflateFromDataIdentity(o); err != nil {
							if errors.Is(err, toml.Error{}) {
								err = nil
							} else {
								err = errors.Wrap(err)
								return
							}
						}

						if err = w(te); err != nil {
							err = errors.Wrap(err)
							return
						}

						return
					},
				); err != nil {
					err = errors.Wrapf(
						err,
						"Transaktion: %s/%s: %s",
						t.Time.Kopf(),
						t.Time.Schwanz(),
						t.Time,
					)

					return
				}

				return
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s konfigStore) ReadOne(
	_ *kennung.Konfig,
) (tt *erworben.Transacted, err error) {
	tt = &erworben.Transacted{
		Sku:  s.StoreUtil.GetKonfig().Sku,
		Akte: s.StoreUtil.GetKonfig().Akte,
	}

	if !tt.Sku.GetTai().IsEmpty() {
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

			defer errors.DeferredCloser(&err, r)

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

			if r, err = s.AkteReader(
				tt.GetAkteSha(),
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

			var sh schnittstellen.ShaLike

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

	if te, err = s.InflateFromDataIdentity(sk); err != nil {
		errors.Wrap(err)
		return
	}

	o = te

	if err = s.updateOne(te); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.LogWriter.Updated(te)

	return
}
