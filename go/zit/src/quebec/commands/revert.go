package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type Revert struct {
	Last bool
}

func init() {
	registerCommandWithQuery(
		"revert",
		func(f *flag.FlagSet) CommandWithQuery {
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
		// gattung.Bestandsaufnahme,
		genres.Repo,
	)
}

func (c Revert) RunWithQuery(
	u *env.Local,
	ms *query.Group,
) (err error) {
	f := func(rt store.RevertId) (err error) {
		if err = u.GetStore().RevertTo(rt); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	switch {
	case c.Last:
		err = c.runRevertFromLast(u, f)

	default:
		err = c.runRevertFromQuery(u, ms, f)
	}

	if err != nil {
		return
	}

	return
}

func (c Revert) runRevertFromQuery(
	u *env.Local,
	eq *query.Group,
	f interfaces.FuncIter[store.RevertId],
) (err error) {
	if err = u.GetStore().QueryTransacted(
		eq,
		func(z *sku.Transacted) (err error) {
			return f(store.RevertId{
				ObjectId: z.GetObjectId(),
				Tai:      z.Metadata.Cache.ParentTai,
			})
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c Revert) runRevertFromLast(
	u *env.Local,
	f interfaces.FuncIter[store.RevertId],
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

			return f(store.RevertId{
				ObjectId: cachedSku.GetObjectId(),
				Tai:      cachedSku.Metadata.Cache.ParentTai,
			})
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
