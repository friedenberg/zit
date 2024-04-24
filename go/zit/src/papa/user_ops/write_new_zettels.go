package user_ops

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/src/november/umwelt"
)

type WriteNewZettels struct {
	*umwelt.Umwelt
	CheckOut        bool
	CheckoutOptions checkout_options.Options
}

func (c WriteNewZettels) RunMany(
	z zettel.ProtoZettel,
	count int,
) (results sku.CheckedOutMutableSet, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	results = collections_value.MakeMutableValueSet[*sku.CheckedOut](
		nil,
	)

	// TODO-P4 modify this to be run once
	for i := 0; i < count; i++ {
		var cz *sku.CheckedOut

		if cz, err = c.runOneAlreadyLocked(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		results.Add(cz)
	}

	return
}

func (c WriteNewZettels) RunOne(
	z zettel.ProtoZettel,
) (result *sku.CheckedOut, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	return c.runOneAlreadyLocked(z)
}

func (c WriteNewZettels) runOneAlreadyLocked(
	pz zettel.ProtoZettel,
) (result *sku.CheckedOut, err error) {
	z := pz.Make()
	defer metadatei.GetPool().Put(z)

	var zt *sku.Transacted

	if zt, err = c.GetStore().Create(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	result = &sku.CheckedOut{}

	if err = result.Internal.SetFromSkuLike(zt); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.CheckOut {
		// TODO-P4 separate creation and checkout into two ops to allow for
		// optimistic
		// unlocking
		if result, err = c.GetStore().CheckoutOne(
			c.CheckoutOptions,
			&result.Internal,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
