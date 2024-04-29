package commands

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/juliett/query"
	"code.linenisgreat.com/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/src/papa/user_ops"
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
	qg *query.Group,
) (err error) {
	opCheckout := user_ops.Checkout{
		Umwelt:  u,
		Options: c.CheckoutOptions,
	}

	if _, err = opCheckout.RunQuery(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
