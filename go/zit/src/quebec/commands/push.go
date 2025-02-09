package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
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

func (c Push) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
}

func (cmd Push) Run(req command.Request) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)

	remoteArg := req.Args()[0]
	remote := cmd.MakeArchive(req, remoteArg, localWorkingCopy)

	repoType := remote.GetImmutableConfig().ImmutableConfig.GetRepoType()

	switch repoType {
	case repo_type.TypeWorkingCopy:
		queryGroup := cmd.MakeQueryGroup(
			req,
			query.BuilderOptions(
				query.BuilderOptionsOld(cmd),
				query.BuilderOptionDefaultSigil(
					ids.SigilHistory,
					ids.SigilHidden,
				),
			),
			localWorkingCopy,
			req.Args()[1:],
		)

		if err := remote.(repo.WorkingCopy).PullQueryGroupFromRemote(
			localWorkingCopy,
			queryGroup,
			cmd.WithPrintCopies(true),
		); err != nil {
			localWorkingCopy.CancelWithError(err)
		}

	case repo_type.TypeArchive:
		if args := req.Args()[1:]; len(args) > 0 {
			ui.Err().Printf("remote is archive, ignore arguments: %q", req.Args()[1:])
		}

		cmd.PushAllToArchive(req, localWorkingCopy, remote)

	default:
		req.CancelWithBadRequestf("unsupported repo type: %q", repoType)
	}
}
