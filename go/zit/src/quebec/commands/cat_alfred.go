package commands

import (
	"bufio"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/alfred"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CatAlfred struct {
	genres.Genre
	Command
}

func init() {
	registerCommandWithQuery(
		"cat-alfred",
		func(f *flag.FlagSet) CommandWithQuery {
			c := &CatAlfred{}

			f.Var(&c.Genre, "genre", "extract this element from all matching objects")

			return c
		},
	)
}

func (c CatAlfred) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Tag,
		genres.Type,
		genres.Zettel,
	)
}

func (c CatAlfred) DefaultGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Tag,
		genres.Type,
		genres.Zettel,
	)
}

func (c CatAlfred) RunWithQuery(
	u *env.Env,
	qg *query.Group,
) (err error) {
	// this command does its own error handling
	wo := bufio.NewWriter(u.Out())
	defer errors.DeferredFlusher(&err, wo)

	var aw *alfred.Writer

	if aw, err = alfred.New(
		wo,
		u.GetStore().GetAbbrStore().GetAbbr(),
		u.SkuFormatBoxNoColor(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if err = u.GetStore().QueryTransacted(
		qg,
		func(object *sku.Transacted) (err error) {
			switch c.Genre {
			case genres.Tag:
				for t := range object.Metadata.Tags.All() {
					var tagObject *sku.Transacted

					if tagObject, err = u.GetStore().ReadTransactedFromObjectId(
						t,
					); err != nil {
						err = errors.Wrap(err)
						return
					}

					if err = aw.PrintOne(tagObject); err != nil {
						err = errors.Wrap(err)
						return
					}
				}
			case genres.Type:
				tipe := object.GetType()

				if tipe.GetType().IsEmpty() {
					return
				}

				if object, err = u.GetStore().ReadTransactedFromObjectId(
					tipe.GetType(),
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = aw.PrintOne(object); err != nil {
					err = errors.Wrap(err)
					return
				}

			default:
				if err = aw.PrintOne(object); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	); err != nil {
		aw.WriteError(err)
		return
	}

	return
}
