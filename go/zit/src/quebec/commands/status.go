package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Status struct{}

func init() {
	registerCommandWithQuery(
		"status",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Status{}

			return c
		},
	)
}

func (c Status) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Status) ModifyBuilder(
	b *query.Builder,
) {
	b.WithHidden(nil)
}

func (c Status) RunWithQuery(
	u *env.Env,
	eqwk *query.Group,
) (err error) {
	pcol := u.PrinterCheckedOutForKasten(eqwk.RepoId)

	if err = u.GetStore().QueryCheckedOut(
		eqwk,
		func(co sku.CheckedOutLike) (err error) {
			if err = pcol(co); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// if err = u.GetStore().QueryUnsure(
	// 	eqwk,
	// 	sku.UnsureMatchOptions{
	// 		UnsureMatchType: sku.UnsureMatchTypeMetadataWithoutTaiHistory | sku.UnsureMatchTypeDescription,
	// 	},
	// 	u.PrinterMatching(),
	// ); err != nil {
	// 	err = errors.Wrap(err)
	// 	return
	// }

	return
}
