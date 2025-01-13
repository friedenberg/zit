package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Checkin struct {
	IgnoreBlob bool
	Proto      sku.Proto

	command_components.Checkout

	CheckoutBlobAndRun string
	OpenBlob           bool
}

func init() {
	f := func(f *flag.FlagSet) CommandWithQuery {
		c := &Checkin{}

		f.BoolVar(
			&c.IgnoreBlob,
			"ignore-blob",
			false,
			"do not change the blob",
		)

		f.StringVar(
			&c.CheckoutBlobAndRun,
			"each-blob",
			"",
			"checkout each Blob and run a utility",
		)

		c.Proto.SetFlagSet(f)
		c.Checkout.SetFlagSet(f)

		return c
	}

	registerCommandWithQuery("checkin", f)
	registerCommandWithQuery("add", f)
	registerCommandWithQuery("save", f)
}

func (c Checkin) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c *Checkin) ModifyBuilder(b *query.Builder) {
	b.
		WithDefaultSigil(ids.SigilExternal).
		WithRequireNonEmptyQuery()
}

func (c Checkin) RunWithQuery(u *read_write_repo_local.Repo, qg *query.Group) {
	op := user_ops.Checkin{
		Delete:             c.Delete,
		Organize:           c.Organize,
		Proto:              c.Proto,
		CheckoutBlobAndRun: c.CheckoutBlobAndRun,
		OpenBlob:           c.OpenBlob,
	}

	// TODO add auto dot operator
	if err := op.Run(u, qg); err != nil {
		u.CancelWithError(err)
	}
}
