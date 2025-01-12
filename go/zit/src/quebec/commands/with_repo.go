package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type commandWithRepo struct {
	*flag.FlagSet
	command_components.RepoLocal
	Command CommandWithRepo
}

func (cmd commandWithRepo) GetFlagSet() *flag.FlagSet {
	return cmd.FlagSet
}

func (cmd commandWithRepo) Run(
	dependencies Dependencies,
) {
	var envOptions env.Options

	if og, ok := cmd.Command.(env.OptionsGetter); ok {
		envOptions = og.GetEnvOptions()
	}

	repoOptions := read_write_repo_local.OptionsEmpty

	if og, ok := cmd.Command.(read_write_repo_local.OptionsGetter); ok {
		repoOptions = og.GetLocalRepoOptions()
	}

	cmdArgs := cmd.Args()

	repo := cmd.MakeRepoLocal(
		dependencies.Context,
		dependencies.Config,
		envOptions,
		repoOptions,
	)

	switch {
	case repo.GetConfig().Complete:
		t := cmd.getCommandCompletionWithRepo(dependencies)
		t.CompleteWithRepo(repo, cmdArgs...)

	default:
		cmd.Command.RunWithRepo(repo, cmdArgs...)
	}
}

func (cmd commandWithRepo) getCommandCompletionWithRepo(
	dependencies Dependencies,
) (t CommandCompletionWithRepo) {
	haystack := any(cmd.Command)

	for {
		switch c := haystack.(type) {
		case *commandWithQuery:
			t = c
			return

		case CommandCompletionWithRepo:
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
