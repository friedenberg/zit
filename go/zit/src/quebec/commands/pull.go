package commands

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Pull struct{}

func init() {
	registerCommand(
		"pull",
		&commandWithLocalWorkingCopy{
			Command: &commandWithRemoteAndQuery{
				Command: &Pull{},
			},
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

func (c Pull) Run(
	local *local_working_copy.Repo,
	remote repo.WorkingCopy,
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
