package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Checkin struct {
	Delete     bool
	IgnoreBlob bool
	Proto      sku.Proto

	Organize           bool
	CheckoutBlobAndRun string
	OpenBlob           bool
	Edit               bool
}

func init() {
	f := func(f *flag.FlagSet) CommandWithQuery {
		c := &Checkin{}

		f.BoolVar(&c.Delete, "delete", false, "the checked-out file")

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

		f.BoolVar(
			&c.Edit,
			"edit",
			false,
			"edit each checked in object",
		)

		f.BoolVar(&c.Organize, "organize", false, "")

		c.Proto.AddToFlagSet(f)

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

func (c Checkin) RunWithQuery(
	u *env.Local,
	qg *query.Group,
) (err error) {
	op := user_ops.Checkin{
		Delete:             c.Delete,
		Organize:           c.Organize,
		Proto:              c.Proto,
		CheckoutBlobAndRun: c.CheckoutBlobAndRun,
		OpenBlob:           c.OpenBlob,
	}

	// TODO add auto dot operator
	if err = op.Run(u, qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
