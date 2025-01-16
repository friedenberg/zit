package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type commandWithLocalWorkingCopy struct {
	*flag.FlagSet
	command_components.LocalWorkingCopy
	Command WithLocalWorkingCopy
}

func (cmd *commandWithLocalWorkingCopy) GetCommand2() Command {
	return cmd
}

func (cmd *commandWithLocalWorkingCopy) SetFlagSet(f *flag.FlagSet) {
	cmd.FlagSet = f

	if cmp, ok := cmd.Command.(interfaces.CommandComponent); ok {
		cmp.SetFlagSet(f)
	}

	cmd.LocalWorkingCopy.SetFlagSet(f)
}

func (cmd commandWithLocalWorkingCopy) GetFlagSet() *flag.FlagSet {
	return cmd.FlagSet
}

func (cmd commandWithLocalWorkingCopy) Run(
	dependencies command.Dep,
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

	repo := cmd.MakeLocalWorkingCopy(
		dependencies,
		envOptions,
		repoOptions,
	)

	switch {
	case repo.GetConfig().Complete:
		t := cmd.getCommandCompletionWithRepo(dependencies)
		t.Complete(repo, cmdArgs...)

	default:
		cmd.Command.Run(repo, cmdArgs...)
	}
}

func (cmd commandWithLocalWorkingCopy) getCommandCompletionWithRepo(
	dependencies command.Dep,
) (t CompleteWithRepo) {
	haystack := any(cmd.Command)

	for {
		switch c := haystack.(type) {
		case *commandWithQuery:
			t = c
			return

		case CompleteWithRepo:
			t = c
			return

		default:
			dependencies.CancelWithBadRequestf(
				"Command does not support completion: %T",
				c,
			)
		}
	}
}
