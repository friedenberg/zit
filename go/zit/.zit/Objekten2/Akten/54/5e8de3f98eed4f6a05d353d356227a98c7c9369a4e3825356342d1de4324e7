package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/echo/checked_out_state"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) checkoutOneIfNecessary(
	options checkout_options.Options,
	tg sku.TransactedGetter,
) (co *sku.CheckedOut, i *sku.FSItem, err error) {
	internal := tg.GetSku()
	co = GetCheckedOutPool().Get()

	sku.Resetter.ResetWith(co.GetSku(), internal)

	var alreadyCheckedOut bool

	if i, alreadyCheckedOut, err = s.prepareFSItemForCheckOut(options, co); err != nil {
		err = errors.Wrap(err)
		return
	}

	if alreadyCheckedOut && !s.shouldCheckOut(options, co, true) {
		if err = s.WriteFSItemToExternal(i, co.GetSkuExternal()); err != nil {
			err = errors.Wrap(err)
			return
		}

		co.SetState(checked_out_state.CheckedOut)

		return
	}

	if err = s.checkoutOneForReal(
		options,
		co,
		i,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) prepareFSItemForCheckOut(
	options checkout_options.Options,
	co *sku.CheckedOut,
) (item *sku.FSItem, alreadyCheckedOut bool, err error) {
	if s.config.IsDryRun() || options.Path == checkout_options.PathTempLocal {
		item = &sku.FSItem{}
		item.Reset()
		return
	}

	if item, alreadyCheckedOut = s.Get(co.GetSku().GetObjectId()); alreadyCheckedOut {
		if err = s.HydrateExternalFromItem(
			sku.CommitOptions{
				Mode: object_mode.ModeRealizeSansProto,
			},
			item,
			co.GetSku(),
			co.GetSkuExternal(),
		); err != nil {
			if sku.IsErrMergeConflict(err) && options.AllowConflicted {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}
	} else {
		if item, err = s.ReadFSItemFromExternal(co.GetSkuExternal()); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	// sku.DetermineState(co, true)

	return
}
