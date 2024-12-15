package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Edit struct {
	// TODO-P3 add force
	Delete       bool
	CheckoutMode checkout_mode.Mode
}

func init() {
	registerCommandWithQuery(
		"edit",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Edit{
				CheckoutMode: checkout_mode.MetadataOnly,
			}

			f.BoolVar(
				&c.Delete,
				"delete",
				false,
				"delete the zettel and blob after successful checkin",
			)
			f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")

			return c
		},
	)
}

func (c Edit) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Tag,
		genres.Zettel,
		genres.Type,
		genres.Repo,
	)
}

func (c Edit) DefaultGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Tag,
		genres.Zettel,
		genres.Type,
		genres.Repo,
	)
}

func (c Edit) RunWithQuery(
	u *env.Local,
	eqwk *query.Group,
) (err error) {
	options := checkout_options.Options{
		CheckoutMode: c.CheckoutMode,
	}

	opEdit := user_ops.Checkout{
		Local:     u,
		Options: options,
		Edit:    true,
	}

	if _, err = opEdit.RunQuery(eqwk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
