package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Pull struct{}

func init() {
	registerCommandWithRemoteAndQuery(
		"pull",
		func(f *flag.FlagSet) CommandWithRemoteAndQuery {
			c := &Pull{}

			return c
		},
	)
}

func (c Pull) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Pull) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
	// return ids.MakeGenre(genres.TrueGenre()...)
}

func (c Pull) RunWithRemoteAndQuery(
	local *repo_local.Local,
	remote repo.Repo,
	qg *query.Group,
) {
	if err := local.PullQueryGroupFromRemote(
		remote,
		qg,
		true,
	); err != nil {
		local.CancelWithError(err)
		return
	}
}
