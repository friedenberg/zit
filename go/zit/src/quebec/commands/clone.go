package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register(
		"clone",
		&Clone{
			Genesis: command_components.Genesis{
				BigBang: repo_layout.BigBang{
					ExcludeDefaultType: true,
				},
			},
		},
	)
}

type Clone struct {
	command_components.Genesis
	command_components.RemoteTransfer
	command_components.QueryGroup
}

func (cmd *Clone) SetFlagSet(f *flag.FlagSet) {
	cmd.Genesis.SetFlagSet(f)
	cmd.RemoteTransfer.SetFlagSet(f)
	cmd.QueryGroup.SetFlagSet(f)

	// must happen after genesis set flag set as cmd.Config is nil until then
	cmd.Config.RepoType = repo_type.TypeWorkingCopy
}

func (c Clone) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Clone) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
	// return ids.MakeGenre(genres.TrueGenre()...)
}

func (cmd Clone) Run(req command.Request) {
	repoGeneric := cmd.OnTheFirstDay(req)

	remote := cmd.MakeRemoteWorkingCopy(
		req,
		req.Args()[0],
	)

	switch local := repoGeneric.(type) {
	default:
		req.CancelWithBadRequestf(
			"unsupported repo type: %q (%T)",
			local.GetRepoType(),
			local,
		)

	case repo.WorkingCopy:
		qg := cmd.MakeQueryGroup(
			req,
			query.MakeBuilderOptions(cmd),
			local,
			req.Args()[1:],
		)

		if err := local.PullQueryGroupFromRemote(
			remote,
			qg,
			cmd.RemoteTransferOptions.WithPrintCopies(true),
		); err != nil {
			req.CancelWithError(err)
		}

	case repo.Repo:
		cmd.PushAllToArchive(req, local, remote)
	}
}
