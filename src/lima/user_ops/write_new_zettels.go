package user_ops

import (
	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/echo/umwelt"
	"github.com/friedenberg/zit/src/foxtrot/zettel"
	zettel_stored "github.com/friedenberg/zit/src/golf/zettel_stored"
	"github.com/friedenberg/zit/src/india/zettel_checked_out"
	store_working_directory "github.com/friedenberg/zit/src/juliett/store_working_directory"
	"github.com/friedenberg/zit/src/kilo/store_with_lock"
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
	var tz zettel_stored.Transacted

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
