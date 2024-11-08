package store_fs

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/object_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (s *Store) prepareFSItemForCheckOut2(
	options checkout_options.Options,
	co *sku.CheckedOut,
) (item *sku.FSItem, err error) {
	if s.config.IsDryRun() {
		item = &sku.FSItem{}
		return
	}

	var ok bool

	if item, ok = s.Get(&co.Internal.ObjectId); !ok {
		if item, err = s.ReadFSItemFromExternal(&co.External); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = s.readIntoExternalFromItem(
		sku.CommitOptions{
			Mode: object_mode.ModeRealizeSansProto,
		},
		item,
		&co.Internal,
		&co.External,
	); err != nil {
		if errors.Is(err, sku.ErrExternalHasConflictMarker) && options.AllowConflicted {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	sku.DetermineState(co, true)

	return
}

func (s *Store) checkoutOneNew2(
	options checkout_options.Options,
	tg sku.TransactedGetter,
) (co *sku.CheckedOut, i *sku.FSItem, err error) {
	internal := tg.GetSku()
	co = GetCheckedOutPool().Get()

	sku.Resetter.ResetWith(&co.Internal, internal)

	if i, err = s.prepareFSItemForCheckOut2(options, co); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !s.shouldCheckOut(options, co, true) {
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
