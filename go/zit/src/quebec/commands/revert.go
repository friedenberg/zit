package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	registerCommand("revert", &Revert{})
}

type Revert struct {
	command_components.LocalWorkingCopyWithQueryGroup

	Last bool
}

func (c *Revert) SetFlagSet(f *flag.FlagSet) {
	c.LocalWorkingCopyWithQueryGroup.SetFlagSet(f)
	f.BoolVar(&c.Last, "last", false, "revert the last changes")
}

func (c Revert) CompletionGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
		genres.InventoryList,
		genres.Repo,
	)
}

func (c Revert) DefaultGenres() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
		genres.Repo,
	)
}

func (cmd Revert) Run(dep command.Dep) {
	localWorkingCopy, queryGroup := cmd.MakeLocalWorkingCopyAndQueryGroup(
		dep,
		query.MakeBuilderOptions(cmd),
	)

	localWorkingCopy.Must(localWorkingCopy.Lock)

	switch {
	case cmd.Last:
		if err := cmd.runRevertFromLast(localWorkingCopy); err != nil {
			localWorkingCopy.CancelWithError(err)
		}

	default:
		if err := cmd.runRevertFromQuery(localWorkingCopy, queryGroup); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	localWorkingCopy.Must(localWorkingCopy.Unlock)
}

func (c Revert) runRevertFromQuery(
	u *local_working_copy.Repo,
	eq *query.Group,
) (err error) {
	if err = u.GetStore().QueryTransacted(
		eq,
		func(z *sku.Transacted) (err error) {
			rt := store.RevertId{
				ObjectId: z.GetObjectId(),
				Tai:      z.Metadata.Cache.ParentTai,
			}

			if err = u.GetStore().RevertTo(rt); err != nil {
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

func (c Revert) runRevertFromLast(
	u *local_working_copy.Repo,
) (err error) {
	s := u.GetStore()

	var b *sku.Transacted

	if b, err = s.GetInventoryListStore().ReadLast(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = s.GetInventoryListStore().StreamInventoryList(
		b.GetBlobSha(),
		func(z *sku.Transacted) (err error) {
			var cachedSku *sku.Transacted

			if cachedSku, err = u.GetStore().GetStreamIndex().ReadOneObjectIdTai(
				z.GetObjectId(),
				z.GetTai(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer sku.GetTransactedPool().Put(cachedSku)

			rt := store.RevertId{
				ObjectId: cachedSku.GetObjectId(),
				Tai:      cachedSku.Metadata.Cache.ParentTai,
			}

			if err = u.GetStore().RevertTo(rt); err != nil {
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
