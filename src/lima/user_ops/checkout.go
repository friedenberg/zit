package user_ops

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/id_set"
	"github.com/friedenberg/zit/src/hotel/zettel_checked_out"
	"github.com/friedenberg/zit/src/india/store_working_directory"
	"github.com/friedenberg/zit/src/kilo/umwelt"
)

type Checkout struct {
	*umwelt.Umwelt
	store_working_directory.CheckoutOptions
}

func (c Checkout) RunMany(
	ids id_set.Set,
) (results zettel_checked_out.Set, err error) {
	if results, err = c.StoreWorkingDirectory().Checkout(c.CheckoutOptions, ids); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
