package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
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
				CheckoutMode: checkout_mode.ModeObjekteOnly,
			}

			f.BoolVar(
				&c.Delete,
				"delete",
				false,
				"delete the zettel and akte after successful checkin",
			)
			f.Var(&c.CheckoutMode, "mode", "mode for checking out the zettel")

			return c
		},
	)
}

func (c Edit) CompletionGattung() ids.Genre {
	return ids.MakeGenre(
		gattung.Etikett,
		gattung.Zettel,
		gattung.Typ,
		gattung.Kasten,
	)
}

func (c Edit) DefaultGattungen() ids.Genre {
	return ids.MakeGenre(
		gattung.Etikett,
		gattung.Zettel,
		gattung.Typ,
		gattung.Kasten,
	)
}

func (c Edit) RunWithQuery(
	u *umwelt.Umwelt,
	eqwk *query.Group,
) (err error) {
	options := checkout_options.Options{
		CheckoutMode: c.CheckoutMode,
	}

	opEdit := user_ops.Checkout{
		Umwelt:  u,
		Options: options,
		Edit:    true,
	}

	if _, err = opEdit.RunQuery(eqwk); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
