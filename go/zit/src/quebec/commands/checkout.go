package commands

import (
	"flag"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
	"code.linenisgreat.com/zit-go/src/charlie/checkout_options"
	"code.linenisgreat.com/zit-go/src/charlie/gattung"
	"code.linenisgreat.com/zit-go/src/delta/gattungen"
	"code.linenisgreat.com/zit-go/src/hotel/sku"
	"code.linenisgreat.com/zit-go/src/india/matcher"
	"code.linenisgreat.com/zit-go/src/oscar/umwelt"
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
