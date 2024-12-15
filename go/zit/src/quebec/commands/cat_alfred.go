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
	u *env.Local,
	qg *query.Group,
) (err error) {
	// this command does its own error handling
	defer func() {
		err = nil
	}()

	wo := bufio.NewWriter(u.Out())
	defer errors.DeferredFlusher(&err, wo)

	var aiw alfred.Writer

	itemPool := alfred.MakeItemPool()

	switch c.Genre {
	case genres.Type, genres.Tag:
		if aiw, err = alfred.NewDebouncingWriter(u.Out()); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		if aiw, err = alfred.NewWriter(u.Out(), itemPool); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	var aw *alfred_sku.Writer

	if aw, err = alfred_sku.New(
		wo,
		u.GetStore().GetAbbrStore().GetAbbr(),
		u.SkuFormatBoxTransactedNoColor(),
		aiw,
		itemPool,
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

	return
}
