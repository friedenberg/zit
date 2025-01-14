package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type commandWithArchive struct {
	*flag.FlagSet
	command_components.Repo
	Command CommandWithArchive
}

func (cmd *commandWithArchive) GetCommand2() Command2 {
	return cmd
}

func (cmd *commandWithArchive) SetFlagSet(f *flag.FlagSet) {
	cmd.FlagSet = f

	if cmp, ok := cmd.Command.(interfaces.CommandComponent); ok {
		cmp.SetFlagSet(f)
	}

	cmd.Repo.SetFlagSet(f)
}

func (cmd commandWithArchive) GetFlagSet() *flag.FlagSet {
	return cmd.FlagSet
}

func (cmd commandWithArchive) Run(
	dependencies Dependencies,
) {
	var envOptions env.Options

	if og, ok := cmd.Command.(env.OptionsGetter); ok {
		envOptions = og.GetEnvOptions()
	}

	repoOptions := local_working_copy.OptionsEmpty

	if og, ok := cmd.Command.(local_working_copy.OptionsGetter); ok {
		repoOptions = og.GetLocalRepoOptions()
	}

	cmdArgs := cmd.Args()

	repo := cmd.MakeArchive(
		dependencies.Context,
		dependencies.Config,
		envOptions,
		repoOptions,
	)

	cmd.Command.RunWithArchive(repo, cmdArgs...)
}
