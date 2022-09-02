package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/umwelt"
	"github.com/friedenberg/zit/src/delta/zettel"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/juliett/store_with_lock"
)

type WriteNewZettels struct {
	*umwelt.Umwelt
	store_working_directory.CheckoutOptions
}

func (c WriteNewZettels) RunMany(
	store store_with_lock.Store,
	z zettel.Zettel,
	count int,
) (results zettel_checked_out.Set, err error) {
	results = zettel_checked_out.MakeSetUnique(count)

	//TODO modify this to be run once
	for i := 0; i < count; i++ {
		var cz zettel_checked_out.Zettel

		if cz, err = c.RunOne(store, z); err != nil {
			err = errors.Wrap(err)
			return
		}

		results.Add(cz)
	}

	return
}

func (c WriteNewZettels) RunOne(
	store store_with_lock.Store,
	z zettel.Zettel,
) (result zettel_checked_out.Zettel, err error) {
	var tz zettel_transacted.Zettel

	if tz, err = store.StoreObjekten().Create(z); err != nil {
		err = errors.Wrap(err)
		return
	}

	if result, err = store.StoreWorkingDirectory().CheckoutOne(c.CheckoutOptions, tz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
