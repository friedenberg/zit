package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
	"code.linenisgreat.com/zit/go/zit/src/papa/remote_push"
)

type Push struct {
	genres.Genre
}

func init() {
	registerCommand(
		"push",
		func(f *flag.FlagSet) Command {
			c := &Push{
				Genre: genres.Zettel,
			}

			f.Var(&c.Genre, "gattung", "Gattung")

			return c
		},
	)
}

func (c Push) CompletionGattung() ids.Genre {
	return ids.MakeGenre(
		genres.Zettel,
		genres.Tag,
		genres.Type,
		genres.InventoryList,
		genres.Repo,
	)
}

func (c Push) Run(u *umwelt.Umwelt, args ...string) (err error) {
	if len(args) == 0 {
		err = errors.Normalf("must specify kasten to push from")
		return
	}

	from := args[0]

	if len(args) > 1 {
		args = args[1:]
	} else {
		err = errors.Normalf("Nothing to push.")
		return
	}

	builder := u.MakeQueryBuilderExcludingHidden(
		ids.MakeGenre(),
	)

	var qg *query.Group

	if qg, err = builder.BuildQueryGroupWithRepoId(
		ids.RepoId{},
		sku.ExternalQueryOptions{},
		args...,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = u.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, u.Unlock)

	var client remote_push.Client

	if client, err = remote_push.MakeClient(u, from); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = client.SendNeededSkus(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, client)

	return
}
