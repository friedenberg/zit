package store_fs

import (
	"path"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/mike/store_util"
)

type Store struct {
	store_util.StoreUtil

	sonnenaufgang kennung.Time
	standort.Standort

	checkedOutLogPrinter schnittstellen.FuncIter[*sku.CheckedOut]
}

func New(
	su store_util.StoreUtil,
	t kennung.Time,
	st standort.Standort,
) (s *Store, err error) {
	s = &Store{
		StoreUtil:     su,
		sonnenaufgang: t,
		Standort:      st,
	}

	return
}

func (s *Store) SetCheckedOutLogPrinter(
	zelw schnittstellen.FuncIter[*sku.CheckedOut],
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

func (s *Store) readOneExternal(
	fs *cwd.CwdFiles,
	sk2 *sku.Transacted,
) (co *sku.CheckedOut, err error) {
	// TODO-P3 pool
	co = &sku.CheckedOut{
		Internal: *sk2,
	}

	ok := false

	var e *sku.ExternalMaybe

	if e, ok = fs.Get(sk2.Kennung); !ok {
		err = iter.MakeErrStopIteration()
		return
	}

	var e2 *sku.External

	if e2, err = s.ReadOneExternal(e, sk2); err != nil {
		if errors.IsNotExist(err) {
			err = iter.MakeErrStopIteration()
		} else {
			err = errors.Wrapf(err, "Cwd: %#v", e)
		}

		return
	}

	co.External = *e2

	return
}

func (s *Store) ReadFiles(
	fs *cwd.CwdFiles,
	fq objekte.FuncReaderTransactedLikePtr,
	f schnittstellen.FuncIter[*sku.CheckedOut],
) (err error) {
	if err = fq(
		iter.MakeChain(
			func(et *sku.Transacted) (err error) {
				var col *sku.CheckedOut

				et1 := sku.GetTransactedPool().Get()

				if err = et1.SetFromSkuLike(et); err != nil {
					err = errors.Wrap(err)
					return
				}

				if col, err = s.readOneExternal(fs, et1); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = col.Internal.SetFromSkuLike(et); err != nil {
					err = errors.Wrap(err)
					return
				}

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

				if err = s.GetAbbrStore().Exists(k); err == nil {
					err = iter.MakeErrStopIteration()
					return
				}

				err = nil

				tco := &sku.CheckedOut{}
				var tcoe *sku.External

				if tcoe, err = s.ReadOneExternal(
					il,
					nil,
				); err != nil {
					if errors.IsNotExist(err) {
						err = iter.MakeErrStopIteration()
					} else {
						err = errors.Wrapf(err, "%#v", il)
					}

					return
				}

				tco.Internal = tcoe.Transacted
				tco.External = *tcoe
				tco.State = checked_out_state.StateUntracked

				if err = f(tco); err != nil {
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

	return
}

// func (s *Store) ReadOneExternal(
// 	em *sku.ExternalMaybe,
// 	t sku.SkuLikePtr,
// ) (e *sku.External, err error) {
// 	var m checkout_mode.Mode

// 	if m, err = em.GetFDs().GetCheckoutMode(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	e = &sku.External{}

// 	if err = e.ResetWithExternalMaybe(*em); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	var t1 sku.SkuLikePtr

// 	if t != nil {
// 		t1 = t
// 	}

// 	switch m {
// 	case checkout_mode.ModeAkteOnly:
// 		if err = s.ReadOneExternalAkte(e, t1); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}

// 	case checkout_mode.ModeObjekteOnly, checkout_mode.ModeObjekteAndAkte:
// 		if err = s.ReadOneExternalObjekte(e, t1); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	return
// }

// func (s *Store) ReadOneExternalObjekte(
// 	e sku.SkuLikeExternalPtr,
// 	t sku.SkuLikePtr,
// ) (err error) {
// 	var f *os.File

// 	if f, err = files.Open(e.GetObjekteFD().Path); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	defer errors.DeferredCloser(&err, f)

// 	if t != nil {
// 		e.GetMetadateiPtr().ResetWith(t.GetMetadatei())
// 	}

// 	if _, err = s.metadateiTextParser.ParseMetadatei(f, e); err != nil {
// 		err = errors.Wrapf(err, "%s", f.Name())
// 		return
// 	}

// 	return
// }

// func (s *Store) ReadOneExternalAkte(
// 	e sku.SkuLikeExternalPtr,
// 	t sku.SkuLikePtr,
// ) (err error) {
// 	e.SetMetadatei(t.GetMetadatei())

// 	var aw sha.WriteCloser

// 	if aw, err = s.AkteWriter(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	defer errors.DeferredCloser(&err, aw)

// 	var f *os.File

// 	if f, err = files.OpenExclusiveReadOnly(
// 		e.GetAkteFD().Path,
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	defer errors.DeferredCloser(&err, f)

// 	if _, err = io.Copy(aw, f); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	sh := sha.Make(aw.GetShaLike())
// 	e.GetMetadateiPtr().AkteSha = sh

// 	return
// }
