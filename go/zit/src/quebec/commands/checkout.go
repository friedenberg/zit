package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Checkout struct {
	CheckoutOptions checkout_options.Options
}

func init() {
	registerCommandWithQuery(
		"checkout",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Checkout{
				CheckoutOptions: checkout_options.Options{
					CheckoutMode: checkout_mode.MetadataOnly,
				},
			}

			c.CheckoutOptions.AddToFlagSet(f)

			return c
		},
	)
}

func (c Checkout) DefaultGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
	)
}

func (c Checkout) ModifyBuilder(b *query.Builder) {
	b.
		WithPermittedSigil(ids.SigilLatest).
		WithPermittedSigil(ids.SigilHidden).
		WithDefaultGenres(ids.MakeGenre(genres.Zettel)).
		WithRequireNonEmptyQuery()
}

func (c Checkout) RunWithQuery(
	u *env.Env,
	qg *query.Group,
) (err error) {
	opCheckout := user_ops.Checkout{
		Env:     u,
		Options: c.CheckoutOptions,
	}

	if _, err = opCheckout.RunQuery(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
