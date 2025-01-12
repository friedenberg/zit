package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
)

type Pull struct{}

func init() {
	registerCommandWithRemoteAndQuery(
		"pull",
		&Pull{},
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
	local *read_write_repo_local.Repo,
	remote repo.ReadWrite,
	qg *query.Group,
	options repo.RemoteTransferOptions,
) {
	if err := local.PullQueryGroupFromRemote(
		remote,
		qg,
		options.WithPrintCopies(true),
	); err != nil {
		local.CancelWithError(err)
	}
}
