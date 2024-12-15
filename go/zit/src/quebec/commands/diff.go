package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
	"code.linenisgreat.com/zit/go/zit/src/papa/user_ops"
)

type Diff struct{}

// TODO switch to registerCommandWithExternalQuery
func init() {
	registerCommandWithQuery(
		"diff",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Diff{}

			return c
		},
	)
}

func (c Diff) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Diff) ModifyBuilder(
	b *query.Builder,
) {
	b.WithHidden(nil)
}

func (c Diff) RunWithQuery(
	u *env.Local,
	qg *query.Group,
) (err error) {
	o := checkout_options.TextFormatterOptions{
		DoNotWriteEmptyDescription: true,
	}

	opDiffFS := user_ops.Diff{
		Local: u,
		TextFormatterFamily: object_metadata.MakeTextFormatterFamily(
			u.GetDirectoryLayout(),
			nil,
		),
	}

	if err = u.GetStore().QuerySkuType(
		qg,
		func(co sku.SkuType) (err error) {
			if err = opDiffFS.Run(co, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
