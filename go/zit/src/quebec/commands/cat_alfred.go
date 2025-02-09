package commands

import (
	"bufio"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/alfred"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/alfred_sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("cat-alfred", &CatAlfred{})
}

type CatAlfred struct {
	command_components.LocalWorkingCopyWithQueryGroup

	genres.Genre
}

func (cmd *CatAlfred) SetFlagSet(f *flag.FlagSet) {
	cmd.LocalWorkingCopyWithQueryGroup.SetFlagSet(f)
	f.Var(&cmd.Genre, "genre", "extract this element from all matching objects")
}

func (c CatAlfred) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Tag,
		genres.Type,
		genres.Zettel,
	)
}

func (cmd CatAlfred) Run(dep command.Request) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.BuilderOptionsOld(
			cmd,
			query.BuilderOptionDefaultGenres(
				genres.Tag,
				genres.Type,
				genres.Zettel,
			),
		),
	)

	// this command does its own error handling
	wo := bufio.NewWriter(localWorkingCopy.GetUIFile())
	defer localWorkingCopy.MustFlush(wo)

	var aiw alfred.Writer

	itemPool := alfred.MakeItemPool()

	switch cmd.Genre {
	case genres.Type, genres.Tag:
		{
			var err error

			if aiw, err = alfred.NewDebouncingWriter(localWorkingCopy.GetUIFile()); err != nil {
				localWorkingCopy.CancelWithError(err)
			}
		}

	default:
		{
			var err error

			if aiw, err = alfred.NewWriter(localWorkingCopy.GetUIFile(), itemPool); err != nil {
				localWorkingCopy.CancelWithError(err)
			}
		}
	}

	var aw *alfred_sku.Writer

	{
		var err error

		if aw, err = alfred_sku.New(
			wo,
			localWorkingCopy.GetStore().GetAbbrStore().GetAbbr(),
			localWorkingCopy.SkuFormatBoxTransactedNoColor(),
			aiw,
			itemPool,
		); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	defer localWorkingCopy.MustClose(aw)

	if err := localWorkingCopy.GetStore().QueryTransacted(
		queryGroup,
		func(object *sku.Transacted) (err error) {
			switch cmd.Genre {
			case genres.Tag:
				for t := range object.Metadata.GetTags().All() {
					var tagObject *sku.Transacted

					if tagObject, err = localWorkingCopy.GetStore().ReadTransactedFromObjectId(
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

				if object, err = localWorkingCopy.GetStore().ReadTransactedFromObjectId(
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
