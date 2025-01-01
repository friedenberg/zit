package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type Push struct{}

func init() {
	registerCommandWithRemoteAndQuery(
		"push",
		&Push{},
	)
}

func (c Push) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Push) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
}

func (c Push) RunWithRemoteAndQuery(
	local *repo_local.Repo,
	remote repo.Repo,
	qg *query.Group,
	options repo.RemoteTransferOptions,
) {
	if err := remote.PullQueryGroupFromRemote(
		local,
		qg,
		options.WithPrintCopies(true),
	); err != nil {
		local.CancelWithError(err)
	}
}
