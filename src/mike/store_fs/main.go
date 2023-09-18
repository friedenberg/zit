package store_fs

import (
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/delta/etikett_akte"
	"github.com/friedenberg/zit/src/delta/kasten_akte"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/erworben"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/checked_out"
	"github.com/friedenberg/zit/src/kilo/konfig"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/lima/objekte_store"
	"github.com/friedenberg/zit/src/lima/store_objekten"
)

type Store struct {
	sonnenaufgang kennung.Time
	erworben      konfig.Compiled
	standort.Standort

	storeObjekten *store_objekten.Store

	checkedOutLogPrinter schnittstellen.FuncIter[objekte.CheckedOutLikePtr]
}

func New(
	t kennung.Time,
	k konfig.Compiled,
	st standort.Standort,
	storeObjekten *store_objekten.Store,
) (s *Store, err error) {
	s = &Store{
		sonnenaufgang: t,
		erworben:      k,
		Standort:      st,
		storeObjekten: storeObjekten,
	}

	return
}

func (s *Store) SetCheckedOutLogPrinter(
	zelw schnittstellen.FuncIter[objekte.CheckedOutLikePtr],
) {
	s.checkedOutLogPrinter = zelw
}

// TODO-P3 move to standort
func (s Store) IndexFilePath() string {
	return path.Join(s.Cwd(), ".ZitCheckoutStoreIndex")
}

func (s Store) Flush() (err error) {
	return
}

func (s *Store) ReadFiles(
	fs *cwd.CwdFiles,
	ms matcher.Query,
	f schnittstellen.FuncIter[objekte.CheckedOutLikePtr],
) (err error) {
	zettelEMGR := objekte_store.MakeExternalMaybeGetterReader[
		zettel.Objekte,
		*zettel.Objekte,
		kennung.Hinweis,
		*kennung.Hinweis,
	](
		fs.GetZettel,
		s.storeObjekten.Zettel(),
	)

	etikettEMGR := objekte_store.MakeExternalMaybeGetterReader[
		etikett_akte.V0,
		*etikett_akte.V0,
		kennung.Etikett,
		*kennung.Etikett,
	](
		fs.GetEtikett,
		s.storeObjekten.Etikett(),
	)

	emgr := objekte_store.MakeExternalMaybeGetterReader2(
		fs.Get,
		s.storeObjekten.GetExternalReader2(),
	)

	kastenEMGR := objekte_store.MakeExternalMaybeGetterReader[
		kasten_akte.V0,
		*kasten_akte.V0,
		kennung.Kasten,
		*kennung.Kasten,
	](
		fs.GetKasten,
		s.storeObjekten.Kasten(),
	)

	if err = s.storeObjekten.Query(
		ms,
		iter.MakeChain(
			func(e sku.SkuLikePtr) (err error) {
				log.Log().Printf("trying to read: %s", e.GetSkuLike())
				var col objekte.CheckedOutLikePtr

				switch et := e.(type) {
				case *transacted.Zettel:
					if col, err = zettelEMGR.ReadOne(*et); err != nil {
						var errAkte store_objekten.ErrExternalAkteExtensionMismatch

						if errors.As(err, &errAkte) {
							fs.MarkUnsureAkten(errAkte.Actual)
							log.Log().Printf("unsure akten: %s", et.GetSkuLike())
							err = nil
						} else {
							err = errors.Wrap(err)
						}

						return
					}

				case *transacted.Typ:
					et1 := &sku.Transacted2{}

					if err = et1.SetFromSkuLike(et); err != nil {
						err = errors.Wrap(err)
						return
					}

					if col, err = emgr.ReadOne(et1); err != nil {
						err = errors.Wrap(err)
						return
					}

				// case transacted.Typ:
				// 	et1 := &sku.Transacted2{}

				// 	if err = et1.SetFromSkuLike(et); err != nil {
				// 		err = errors.Wrap(err)
				// 		return
				// 	}

				// 	if col, err = emgr.ReadOne(et1); err != nil {
				// 		err = errors.Wrap(err)
				// 		return
				// 	}

				case *sku.Transacted2:

					if col, err = emgr.ReadOne(et); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *transacted.Kasten:
					if col, err = kastenEMGR.ReadOne(*et); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *transacted.Etikett:
					if col, err = etikettEMGR.ReadOne(*et); err != nil {
						err = errors.Wrap(err)
						return
					}

				case *erworben.Transacted:
					errors.TodoP1("implement checked out konfig?")
					return

				default:
					err = errors.Implement()
					return
				}

				log.Log().Printf("read: %s", col)

				col.DetermineState(false)

				if err = f(col); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = fs.EachCreatableMatchable(
		iter.MakeChain(
			func(il *sku.ExternalMaybe) (err error) {
				k := il.GetKennungLike()

				if err = s.storeObjekten.GetAbbrStore().Exists(k); err == nil {
					err = iter.MakeErrStopIteration()
					return
				}

				err = nil

				switch k.GetGattung() {
				case gattung.Kasten:
					var tco checked_out.Kasten

					if tco.External, err = s.storeObjekten.Kasten().ReadOneExternal(
						*il,
						nil,
					); err != nil {
						if errors.IsNotExist(err) {
							err = iter.MakeErrStopIteration()
						} else {
							err = errors.Wrapf(err, "CwdEtikett: %#v", il)
						}

						return
					}

					tco.State = checked_out_state.StateUntracked

					if err = f(&tco); err != nil {
						err = errors.Wrap(err)
						return
					}

				case gattung.Typ:
					var tco objekte.CheckedOut2

					var e1 *sku.External2

					if e1, err = s.storeObjekten.ReadOneExternal(
						il,
						nil,
					); err != nil {
						if errors.IsNotExist(err) {
							err = iter.MakeErrStopIteration()
						} else {
							err = errors.Wrapf(err, "CwdEtikett: %#v", il)
						}

						return
					}

					tco.External = *e1
					tco.Internal.Kennung = tco.External.Kennung

					tco.State = checked_out_state.StateUntracked

					if err = f(&tco); err != nil {
						err = errors.Wrap(err)
						return
					}

				case gattung.Etikett:
					var tco checked_out.Etikett

					if tco.External, err = s.storeObjekten.Etikett().ReadOneExternal(
						*il,
						nil,
					); err != nil {
						if errors.IsNotExist(err) {
							err = iter.MakeErrStopIteration()
						} else {
							err = errors.Wrapf(err, "CwdEtikett: %#v", il)
						}

						return
					}

					tco.State = checked_out_state.StateUntracked

					if err = f(&tco); err != nil {
						err = errors.Wrap(err)
						return
					}

				default:
					err = errors.Errorf("unsupported id like: %T", il)
				}

				return
			},
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
