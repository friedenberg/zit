package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/bravo/checkout_mode"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Edit struct {
	// TODO-P3 add force
	Workspace bool
	command_components.Checkout
	CheckoutMode checkout_mode.Mode
}

func init() {
	registerCommandWithQuery(
		"edit",
		func(f *flag.FlagSet) WithQuery {
			c := &Edit{
				Workspace:    true,
				CheckoutMode: checkout_mode.MetadataOnly,
			}

			c.Checkout.SetFlagSet(f)

			f.Var(&c.CheckoutMode, "mode", "mode for checking out the object")
			f.BoolVar(&c.Workspace, "use-workspace", true, "checkout the object into the current workspace (CWD)")

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

func (c Edit) RunWithQuery(u *local_working_copy.Repo, eqwk *query.Group) {
	options := checkout_options.Options{
		CheckoutMode: c.CheckoutMode,
	}

	opEdit := user_ops.Checkout{
		Repo:    u,
		Options: options,
		Edit:    true,
	}

	opEdit.Options.Workspace = c.Workspace

	if _, err := opEdit.RunQuery(eqwk); err != nil {
		u.CancelWithError(err)
	}
}
