package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Revert struct {
	Last bool
}

func init() {
	registerCommandWithQuery(
		"revert",
		func(f *flag.FlagSet) WithQuery {
			c := &Revert{}

			f.BoolVar(&c.Last, "last", false, "revert the last changes")

			return c
		},
	)
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

func (c Revert) RunWithQuery(u *local_working_copy.Repo, ms *query.Group) {
	u.Must(u.Lock)

	switch {
	case c.Last:
		if err := c.runRevertFromLast(u); err != nil {
			u.CancelWithError(err)
		}

	default:
		if err := c.runRevertFromQuery(u, ms); err != nil {
			u.CancelWithError(err)
		}
	}

	u.Must(u.Unlock)
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
