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

type Push struct {
	command_components.Repo
	command_components.RemoteTransfer
	command_components.QueryGroup
}

func init() {
	registerCommand(
		"push",
		&Push{},
	)
}

func (cmd *Push) SetFlagSet(f *flag.FlagSet) {
	cmd.RemoteTransfer.SetFlagSet(f)
	cmd.QueryGroup.SetFlagSet(f)
	cmd.Repo.SetFlagSet(f)
}

func (c Push) DefaultSigil() ids.Sigil {
	return ids.MakeSigil(ids.SigilHistory, ids.SigilHidden)
}

func (c Push) DefaultGenres() ids.Genre {
	return ids.MakeGenre(genres.InventoryList)
}

func (cmd Push) Run(
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
		if err := remote.PullQueryGroupFromRemote(
			local,
			qg,
			cmd.WithPrintCopies(true),
		); err != nil {
			local.CancelWithError(err)
		}

	case repo.Archive:
		cmd.PushAllToArchive(local, remote)
	}
}
