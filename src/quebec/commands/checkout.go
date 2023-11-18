package commands

import (
	"flag"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/checkout_options"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/delta/gattungen"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type Checkout struct {
	CheckoutOptions checkout_options.Options
}

func init() {
	registerCommandWithQuery(
		"checkout",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Checkout{}

			c.CheckoutOptions.AddToFlagSet(f)

			return c
		},
	)
}

func (c Checkout) DefaultGattungen() gattungen.Set {
	return gattungen.MakeSet(
		gattung.Zettel,
	)
}

func (c Checkout) RunWithQuery(
	u *umwelt.Umwelt,
	ms matcher.Query,
) (err error) {
	if err = u.StoreObjekten().CheckoutQuery(
		c.CheckoutOptions,
		matcher.MakeFuncReaderTransactedLikePtr(ms, u.StoreObjekten().QueryWithoutCwd),
		func(co *sku.CheckedOut) (err error) {
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
