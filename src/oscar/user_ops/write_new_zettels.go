package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/juliett/zettel"
	"github.com/friedenberg/zit/src/lima/zettel_checked_out"
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
) (results zettel_checked_out.MutableSet, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer c.Unlock()

	results = zettel_checked_out.MakeMutableSetUnique(count)

	//TODO modify this to be run once
	for i := 0; i < count; i++ {
		var cz zettel_checked_out.Zettel

		if cz, err = c.runOneAlreadyLocked(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		results.Add(&cz)
	}

	return
}

func (c WriteNewZettels) RunOne(
	z zettel.ProtoZettel,
) (result zettel_checked_out.Zettel, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer c.Unlock()

	return c.runOneAlreadyLocked(z)
}

func (c WriteNewZettels) runOneAlreadyLocked(
	pz zettel.ProtoZettel,
) (result zettel_checked_out.Zettel, err error) {
	z := pz.Make()

	var zt *zettel.Transacted

	if zt, err = c.StoreObjekten().Zettel().Create(*z); err != nil {
		err = errors.Wrap(err)
		return
	}

	result.Internal = *zt

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