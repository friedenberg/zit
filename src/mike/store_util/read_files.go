package store_util

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/delta/checked_out_state"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/kilo/cwd"
)

func (s *common) ReadOneExternalFS(
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

func (s *common) ReadFiles(
	fs *cwd.CwdFiles,
	fq matcher.FuncReaderTransactedLikePtr,
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

				if col, err = s.ReadOneExternalFS(fs, et1); err != nil {
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
				if err = s.GetAbbrStore().Exists(&il.Kennung); err == nil {
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

				tco.Internal.Kennung = tcoe.Kennung
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
