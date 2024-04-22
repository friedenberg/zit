package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/juliett/query"
	"code.linenisgreat.com/zit/src/oscar/umwelt"
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

func (c Checkout) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung(
		gattung.Zettel,
	)
}

func (c Checkout) RunWithQuery(
	u *umwelt.Umwelt,
	ms *query.Group,
) (err error) {
	if err = u.GetStore().CheckoutQuery(
		c.CheckoutOptions,
		ms,
		func(co *sku.CheckedOut) (err error) {
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
