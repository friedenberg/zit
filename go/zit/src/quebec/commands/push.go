package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
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
	command_components.Query
}

func (cmd *Push) SetFlagSet(flagSet *flag.FlagSet) {
	cmd.RemoteTransfer.SetFlagSet(flagSet)
	cmd.Query.SetFlagSet(flagSet)
	cmd.LocalWorkingCopy.SetFlagSet(flagSet)
}

func (cmd Push) Run(req command.Request) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)

	var remoteObject *sku.Transacted

	{
		var err error

		if remoteObject, err = localWorkingCopy.GetObjectFromObjectId(
			req.PopArg("repo-id"),
		); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	remote := cmd.MakeRemote(req, localWorkingCopy, remoteObject)

	repoType := remote.GetImmutableConfigPublic().ImmutableConfig.GetRepoType()

	switch repoType {
	case repo_type.TypeWorkingCopy:
		queryGroup := cmd.MakeQueryIncludingWorkspace(
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
