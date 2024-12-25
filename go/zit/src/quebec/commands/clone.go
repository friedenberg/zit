package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Clone struct {
	repo_local.BigBang
}

func init() {
	registerCommandWithRemoteAndQueryAndWithoutEnvironment(
		"clone",
		func(f *flag.FlagSet) CommandWithRemoteAndQuery {
			c := &Clone{
				BigBang: repo_local.BigBang{
					Config:             immutable_config.Default(),
					ExcludeDefaultType: true,
				},
			}

			c.BigBang.AddToFlagSet(f)

			return c
		},
	)
}

func (c Clone) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Clone) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
	// return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Clone) RunWithRemoteAndQuery(
	local *repo_local.Local,
	remote repo.Repo,
	qg *query.Group,
) {
	if err := local.Start(c.BigBang); err != nil {
		local.CancelWithError(err)
		return
	}

	if err := local.PullQueryGroupFromRemote(
		remote,
		qg,
		true,
	); err != nil {
		local.CancelWithError(err)
		return
	}

	return
}
