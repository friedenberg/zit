package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type Pull struct {
	command_components.LocalWorkingCopy
	command_components.RemoteTransfer
	command_components.QueryGroup
}

func init() {
	registerCommand(
		"pull",
		&Pull{},
	)
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

func (cmd Pull) Run(
	dependencies Dependencies,
) {
	local := cmd.MakeLocalWorkingCopy(
		dependencies.Context,
		dependencies.Config,
		env.Options{},
		local_working_copy.OptionsEmpty,
	)

	remote := cmd.MakeWorkingCopyFromFlagSet(local.Env, dependencies.FlagSet)

	qg := cmd.MakeQueryGroup(
		query.MakeBuilderOptions(cmd),
		local,
		dependencies.Args()[1:]...,
	)

	switch remote := remote.(type) {
	case repo.WorkingCopy:
		if err := local.PullQueryGroupFromRemote(
			remote,
			qg,
			cmd.WithPrintCopies(true),
		); err != nil {
			local.CancelWithError(err)
		}

	case repo.Archive:
		local.CancelWithBadRequestf(
			"unsupported repo type: %s (%T)",
			remote.GetRepoLayout().GetConfig().GetRepoType(),
			remote,
		)
	}
}
