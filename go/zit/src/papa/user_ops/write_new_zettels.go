package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/checkout_options"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/oscar/umwelt"
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

	if zt, err = c.StoreObjekten().Create(z); err != nil {
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
		if result, err = c.StoreObjekten().CheckoutOne(
			c.CheckoutOptions,
			&result.Internal,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}