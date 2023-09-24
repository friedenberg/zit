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
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/kilo/konfig"
	"github.com/friedenberg/zit/src/lima/cwd"
	"github.com/friedenberg/zit/src/lima/objekte_store"
	"github.com/friedenberg/zit/src/november/store_objekten"
)

type Store struct {
	sonnenaufgang kennung.Time
	erworben      konfig.Compiled
	standort.Standort

	storeObjekten *store_objekten.Store

	checkedOutLogPrinter schnittstellen.FuncIter[*sku.CheckedOut]
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

func (s *Store) ReadFiles(
	fs *cwd.CwdFiles,
	ms matcher.Query,
	f schnittstellen.FuncIter[*sku.CheckedOut],
) (err error) {
	emgr := objekte_store.MakeExternalMaybeGetterReader2(
		fs.Get,
		s.storeObjekten,
	)

	if err = s.storeObjekten.Query(
		ms,
		iter.MakeChain(
			func(et *sku.Transacted) (err error) {
				var col *sku.CheckedOut

				et1 := sku.GetTransactedPool().Get()

				if err = et1.SetFromSkuLike(et); err != nil {
					err = errors.Wrap(err)
					return
				}

				if col, err = emgr.ReadOne(et1); err != nil {
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

				if err = s.storeObjekten.GetAbbrStore().Exists(k); err == nil {
					err = iter.MakeErrStopIteration()
					return
				}

				err = nil

				tco := &sku.CheckedOut{}
				var tcoe *sku.External

				if tcoe, err = s.storeObjekten.ReadOneExternal(
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
