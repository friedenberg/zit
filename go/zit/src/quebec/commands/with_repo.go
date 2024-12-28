package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type commandWithRepo struct {
	Command CommandWithRepo
	*flag.FlagSet
}

func (cmd commandWithRepo) GetFlagSet() *flag.FlagSet {
	return cmd.FlagSet
}

func (cmd commandWithRepo) RunWithDependencies(
	dependencies Dependencies,
) {
	options := repo_local.OptionsEmpty

	if og, ok := cmd.Command.(repo_local.OptionsGetter); ok {
		options = og.GetEnvironmentInitializeOptions()
	}

	cmdArgs := cmd.Args()

	var layout dir_layout.Layout

	{
		var err error

		if layout, err = dir_layout.MakeDefault(
			dependencies.Debug,
		); err != nil {
			dependencies.CancelWithError(err)
		}
	}

	env := env.Make(
		dependencies.Context,
		dependencies.Config,
		layout,
	)

	repo := repo_local.Make(env, options)

	defer dependencies.MustWithContext(repo.GetDirLayout().ResetTempOnExit)
	defer repo.MustFlush(repo)

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
