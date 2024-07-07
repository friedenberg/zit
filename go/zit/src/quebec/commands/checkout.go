package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Checkout struct {
	CheckoutOptions checkout_options.Options
}

func init() {
	registerCommandWithExternalQuery(
		"checkout",
		func(f *flag.FlagSet) CommandWithExternalQuery {
			c := &Checkout{
				CheckoutOptions: checkout_options.Options{
					CheckoutMode: checkout_mode.ModeObjekteOnly,
				},
			}

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

func (c Checkout) RunWithExternalQuery(
	u *umwelt.Umwelt,
	eqwk sku.ExternalQuery,
) (err error) {
	opCheckout := user_ops.Checkout{
		Umwelt:  u,
		Options: c.CheckoutOptions,
	}

	if _, err = opCheckout.RunQuery(eqwk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
