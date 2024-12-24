package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
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
	b.WithHidden(nil).
		WithDefaultSigil(ids.SigilExternal)
}

func (c Status) RunWithQuery(
	u *env.Local,
	qg *query.Group,
) (err error) {
	pcol := u.PrinterCheckedOut(box_format.CheckedOutHeaderState{})

	if err = u.GetStore().QuerySkuType(
		qg,
		func(co sku.SkuType) (err error) {
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

	return
}
