package commands

import (
	"bufio"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/alfred"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/alfred_sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type CatAlfred struct {
	genres.Genre
	CommandWithRepo
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

func (c CatAlfred) RunWithQuery(u *repo_local.Repo, qg *query.Group) {
	// this command does its own error handling
	wo := bufio.NewWriter(u.GetOutFile())
	defer u.MustFlush(wo)

	var aiw alfred.Writer

	itemPool := alfred.MakeItemPool()

	switch c.Genre {
	case genres.Type, genres.Tag:
		{
			var err error

			if aiw, err = alfred.NewDebouncingWriter(u.GetOutFile()); err != nil {
				u.CancelWithError(err)
			}
		}

	default:
		{
			var err error

			if aiw, err = alfred.NewWriter(u.GetOutFile(), itemPool); err != nil {
				u.CancelWithError(err)
			}
		}
	}

	var aw *alfred_sku.Writer

	{
		var err error

		if aw, err = alfred_sku.New(
			wo,
			u.GetStore().GetAbbrStore().GetAbbr(),
			u.SkuFormatBoxTransactedNoColor(),
			aiw,
			itemPool,
		); err != nil {
			u.CancelWithError(err)
		}
	}

	defer u.MustClose(aw)

	if err := u.GetStore().QueryTransacted(
		qg,
		func(object *sku.Transacted) (err error) {
			switch c.Genre {
			case genres.Tag:
				for t := range object.Metadata.GetTags().All() {
					var tagObject *sku.Transacted

					if tagObject, err = u.GetStore().ReadTransactedFromObjectId(
						t,
					); err != nil {
						if collections.IsErrNotFound(err) {
							err = nil
							tagObject = sku.GetTransactedPool().Get()
							defer sku.GetTransactedPool().Put(tagObject)
							tagObject.ObjectId.ResetWithIdLike(t)
						} else {
							err = errors.Wrap(err)
							return
						}
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
		err = nil
		return
	}
}
