package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/juliett/store_working_directory"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type WriteNewZettels struct {
	*umwelt.Umwelt
	CheckOut bool
	store_working_directory.CheckoutOptions
}

func (c WriteNewZettels) RunMany(
	z zettel.Zettel,
	count int,
) (results zettel_checked_out.Set, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer c.Unlock()

	results = zettel_checked_out.MakeSetUnique(count)

	//TODO modify this to be run once
	for i := 0; i < count; i++ {
		var cz zettel_checked_out.Zettel

		if cz, err = c.runOneAlreadyLocked(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		results.Add(cz)
	}

	return
}

func (c WriteNewZettels) RunOne(
	z zettel.Zettel,
) (result zettel_checked_out.Zettel, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer c.Unlock()

	return c.runOneAlreadyLocked(z)
}

func (c WriteNewZettels) runOneAlreadyLocked(
	z zettel.Zettel,
) (result zettel_checked_out.Zettel, err error) {
	if result.Internal, err = c.StoreObjekten().Create(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	if c.CheckOut {
		//TODO separate creation and checkout into two ops to allow for optimistic
		//unlocking
		if result, err = c.StoreWorkingDirectory().CheckoutOne(c.CheckoutOptions, result.Internal); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
