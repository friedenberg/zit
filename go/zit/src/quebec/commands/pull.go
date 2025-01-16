package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
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

func (c Pull) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Pull) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
}

func (cmd Pull) Run(dep command.Request) {
	localWorkingCopy := cmd.MakeLocalWorkingCopy(dep)

	remote := cmd.MakeWorkingCopyFromFlagSet(dep)

	qg := cmd.MakeQueryGroup(
		query.MakeBuilderOptions(cmd),
		localWorkingCopy,
		dep.Args()[1:],
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

	case repo.Archive:
		localWorkingCopy.CancelWithBadRequestf(
			"unsupported repo type: %s (%T)",
			remote.GetRepoLayout().GetConfig().GetRepoType(),
			remote,
		)
	}
}
