package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/juliett/objekte"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/store_util"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type WriteNewZettels struct {
	*umwelt.Umwelt
	CheckOut bool
	store_util.CheckoutOptions
}

func (c WriteNewZettels) RunMany(
	z zettel.ProtoZettel,
	count int,
) (results schnittstellen.MutableSetLike[objekte.CheckedOutLikePtr], err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	results = collections_value.MakeMutableValueSet[objekte.CheckedOutLikePtr](
		nil,
	)

	// TODO-P4 modify this to be run once
	for i := 0; i < count; i++ {
		var cz objekte.CheckedOutLikePtr

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
) (result objekte.CheckedOutLikePtr, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	return c.runOneAlreadyLocked(z)
}

func (c WriteNewZettels) runOneAlreadyLocked(
	pz zettel.ProtoZettel,
) (result objekte.CheckedOutLikePtr, err error) {
	z := pz.Make()

	var zt *transacted.Zettel

	if zt, err = c.StoreObjekten().Zettel().Create(*z); err != nil {
		err = errors.Wrap(err)
		return
	}

	result = &objekte.CheckedOut2{}

	if err = result.GetInternalLikePtr().SetFromSkuLike(zt); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.CheckOut {
		// TODO-P4 separate creation and checkout into two ops to allow for
		// optimistic
		// unlocking
		if result, err = c.StoreWorkingDirectory().CheckoutOneZettel(
			c.CheckoutOptions,
			result.GetInternalLikePtr(),
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
