package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Show struct {
	Format string
}

func init() {
	registerCommandWithQuery(
		"show",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Show{}

			f.StringVar(&c.Format, "format", "log", "format")

			return c
		},
	)
}

func (c Show) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
		genres.InventoryList,
		genres.Repo,
	)
}

func (c Show) DefaultGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
	)
}

func (c Show) RunWithQuery(
	u *env.Env,
	eqwk *query.Group,
) (err error) {
	var f interfaces.FuncIter[*sku.Transacted]

	if f, err = u.MakeFormatFunc(c.Format, u.Out()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.GetStore().QueryWithKasten(
		eqwk,
		iter.MakeSyncSerializer(f),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
