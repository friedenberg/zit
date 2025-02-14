package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("pull", &Pull{})
}

type Pull struct {
	command_components.LocalWorkingCopy
	command_components.RemoteTransfer
	command_components.QueryGroup
}

func (cmd *Pull) SetFlagSet(f *flag.FlagSet) {
	cmd.RemoteTransfer.SetFlagSet(f)
	cmd.QueryGroup.SetFlagSet(f)
	cmd.LocalWorkingCopy.SetFlagSet(f)
}

func (cmd Pull) Run(req command.Request) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(req)

	var object *sku.Transacted

	{
		var err error

		if object, err = localWorkingCopy.GetObjectFromObjectId(
			req.PopArg("repo-id"),
		); err != nil {
			localWorkingCopy.CancelWithError(err)
		}
	}

	remote := cmd.MakeRemote(req, localWorkingCopy, object)

	qg := cmd.MakeQueryGroup(
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

	switch remote := remote.(type) {
	case repo.WorkingCopy:
		if err := localWorkingCopy.PullQueryGroupFromRemote(
			remote,
			qg,
			cmd.WithPrintCopies(true),
		); err != nil {
			localWorkingCopy.CancelWithError(err)
		}

	case repo.Repo:
		localWorkingCopy.CancelWithBadRequestf(
			"unsupported repo type: %s (%T)",
			remote.GetImmutableConfig().ImmutableConfig.GetRepoType(),
			remote,
		)
	}
}
