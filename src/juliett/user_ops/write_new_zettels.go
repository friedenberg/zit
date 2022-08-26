package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	"github.com/friedenberg/zit/src/golf/stored_zettel"
	store_working_directory "github.com/friedenberg/zit/src/hotel/store_working_directory"
	"github.com/friedenberg/zit/src/india/store_with_lock"
	"github.com/friedenberg/zit/src/india/zettel_checked_out"
)

type WriteNewZettels struct {
	*umwelt.Umwelt
	store_working_directory.CheckoutOptions
}

func (c WriteNewZettels) RunMany(
	store store_with_lock.Store,
	zettelen ...zettel.Zettel,
) (results []zettel_checked_out.CheckedOut, err error) {
	results = make([]zettel_checked_out.CheckedOut, 0, len(zettelen))

	for _, z := range zettelen {
		var cz zettel_checked_out.CheckedOut

		if cz, err = c.RunOne(store, z); err != nil {
			err = errors.Error(err)
			return
		}

		results = append(results, cz)
	}

	return
}

func (c WriteNewZettels) RunOne(
	store store_with_lock.Store,
	z zettel.Zettel,
) (result zettel_checked_out.CheckedOut, err error) {
	var tz stored_zettel.Transacted

	if tz, err = store.StoreObjekten().Create(z); err != nil {
		err = errors.Error(err)
		return
	}

	if result, err = store.CheckoutStore().CheckoutOne(c.CheckoutOptions, tz); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
