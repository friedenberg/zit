package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/india/transacted"
	"github.com/friedenberg/zit/src/kilo/zettel"
	"github.com/friedenberg/zit/src/mike/store_fs"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type WriteNewZettels struct {
	*umwelt.Umwelt
	CheckOut bool
	store_fs.CheckoutOptions
}

func (c WriteNewZettels) RunMany(
	z zettel.ProtoZettel,
	count int,
) (results zettel.MutableSetCheckedOut, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	results = zettel.MakeMutableSetCheckedOutUnique(count)

	// TODO-P4 modify this to be run once
	for i := 0; i < count; i++ {
		var cz zettel.CheckedOut

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
) (result zettel.CheckedOut, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	return c.runOneAlreadyLocked(z)
}

func (c WriteNewZettels) runOneAlreadyLocked(
	pz zettel.ProtoZettel,
) (result zettel.CheckedOut, err error) {
	z := pz.Make()

	var zt *transacted.Zettel

	if zt, err = c.StoreObjekten().Zettel().Create(*z); err != nil {
		err = errors.Wrap(err)
		return
	}

	result.Internal = *zt

	if c.CheckOut {
		// TODO-P4 separate creation and checkout into two ops to allow for optimistic
		// unlocking
		if result, err = c.StoreWorkingDirectory().CheckoutOneZettel(
			c.CheckoutOptions,
			result.Internal,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
