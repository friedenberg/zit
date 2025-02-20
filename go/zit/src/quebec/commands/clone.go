package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register(
		"clone",
		&Clone{
			Genesis: command_components.Genesis{
				BigBang: env_repo.BigBang{
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

func (cmd Clone) Run(req command.Request) {
	local := cmd.OnTheFirstDay(req, req.PopArg("new repo id"))

	// TODO offer option to persist remote object, if supported
	remote := cmd.MakeRemote(
		req,
		local,
		cmd.CreateRemoteObject(req, local),
	)

	switch local := local.(type) {
	default:
		req.CancelWithBadRequestf(
			"unsupported repo type: %q (%T)",
			local.GetImmutableConfig().ImmutableConfig.GetRepoType(),
			local,
		)

	case repo.WorkingCopy:
		queryGroup := cmd.MakeQueryGroup(
			req,
			query.BuilderOptionsOld(
				cmd,
				query.BuilderOptionDefaultSigil(
					ids.SigilHistory,
					ids.SigilHidden,
				),
				query.BuilderOptionDefaultGenres(genres.InventoryList),
			),
			local,
			req.PopArgs(),
		)

		if err := local.PullQueryGroupFromRemote(
			remote,
			queryGroup,
			cmd.WithPrintCopies(true),
		); err != nil {
			req.CancelWithError(err)
		}

	case repo.Repo:
		cmd.PushAllToArchive(req, remote, local)
	}
}
