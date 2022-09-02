package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/india/zettel_checked_out"
	"github.com/friedenberg/zit/src/juliett/store_working_directory"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
	"github.com/friedenberg/zit/zettel_transacted"
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
			err = errors.Error(err)
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
	var tz zettel_transacted.Transacted

	if tz, err = store.StoreObjekten().Create(z); err != nil {
		err = errors.Error(err)
		return
	}

	if result, err = store.StoreWorkingDirectory().CheckoutOne(c.CheckoutOptions, tz); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
