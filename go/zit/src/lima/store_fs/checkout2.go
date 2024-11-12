package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
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
		if err = s.WriteFSItemToExternal(i, &co.External); err != nil {
			err = errors.Wrap(err)
			return
		}

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
	if s.config.IsDryRun() {
		item = &sku.FSItem{}
		return
	}

	if item, alreadyCheckedOut = s.Get(&co.GetSku().ObjectId); alreadyCheckedOut {
		if err = s.HydrateExternalFromItem(
			sku.CommitOptions{
				Mode: object_mode.ModeRealizeSansProto,
			},
			item,
			co.GetSku(),
			&co.External,
		); err != nil {
			if errors.Is(err, sku.ErrExternalHasConflictMarker) && options.AllowConflicted {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}
	} else {
		if item, err = s.ReadFSItemFromExternal(&co.External); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	sku.DetermineState(co, true)

	return
}
