package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("push", &Push{})
}

type Push struct {
	command_components.LocalWorkingCopy
	command_components.RemoteTransfer
	command_components.QueryGroup
}

func (cmd *Push) SetFlagSet(f *flag.FlagSet) {
	cmd.RemoteTransfer.SetFlagSet(f)
	cmd.QueryGroup.SetFlagSet(f)
	cmd.LocalWorkingCopy.SetFlagSet(f)
}

func (cmd Push) Run(req command.Request) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)

	remoteArg := req.PopArg("remote arg")
	remote := cmd.MakeArchiveFromArg(req, remoteArg, localWorkingCopy)

	repoType := remote.GetImmutableConfig().ImmutableConfig.GetRepoType()

	switch repoType {
	case repo_type.TypeWorkingCopy:
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
			localWorkingCopy,
			req.PopArgs(),
		)

		if err := remote.(repo.WorkingCopy).PullQueryGroupFromRemote(
			localWorkingCopy,
			queryGroup,
			cmd.WithPrintCopies(true),
		); err != nil {
			localWorkingCopy.CancelWithError(err)
		}

	case repo_type.TypeArchive:
		req.AssertNoMoreArgs()
		cmd.PushAllToArchive(req, localWorkingCopy, remote)

	default:
		req.CancelWithBadRequestf("unsupported repo type: %q", repoType)
	}
}
