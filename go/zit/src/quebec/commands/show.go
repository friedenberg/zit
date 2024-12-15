package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Show struct {
	After  ids.Tai
	Before ids.Tai
	Format string
}

func init() {
	registerCommandWithQuery(
		"show",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &Show{}

			f.StringVar(&c.Format, "format", "log", "format")
			f.Var((*ids.TaiRFC3339Value)(&c.Before), "before", "")
			f.Var((*ids.TaiRFC3339Value)(&c.After), "after", "")

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
	u *env.Local,
	qg *query.Group,
) (err error) {
	var f interfaces.FuncIter[*sku.Transacted]

	if c.Format == "" && qg.IsExactlyOneObjectId() {
		c.Format = "text"
	}

	if f, err = u.MakeFormatFunc(c.Format, u.Out()); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !c.Before.IsEmpty() {
		old := f

		f = func(sk *sku.Transacted) (err error) {
			if !sk.GetTai().Before(c.Before) {
				return
			}

			if err = old(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	if !c.After.IsEmpty() {
		old := f

		f = func(sk *sku.Transacted) (err error) {
			if !sk.GetTai().After(c.After) {
				return
			}

			if err = old(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	if err = u.GetStore().QueryTransacted(
		qg,
		quiter.MakeSyncSerializer(f),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
