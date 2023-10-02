package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/store_util"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type WriteNewZettels struct {
	*umwelt.Umwelt
	CheckOut bool
	store_util.CheckoutOptions
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

	results = collections_ptr.MakeMutableValueSet[sku.CheckedOut, *sku.CheckedOut](
		nil,
	)

	// TODO-P4 modify this to be run once
	for i := 0; i < count; i++ {
		var cz *sku.CheckedOut

		if cz, err = c.runOneAlreadyLocked(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		results.AddPtr(cz)
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

	var zt *sku.Transacted

	if zt, err = c.StoreObjekten().Zettel().Create(*z); err != nil {
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
