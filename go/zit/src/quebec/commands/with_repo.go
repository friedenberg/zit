package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
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
	// TODO use options when making dirLayout
	var dirLayout dir_layout.Layout

	{
		var err error

		if dirLayout, err = dir_layout.MakeDefault(
			dependencies.Debug,
		); err != nil {
			dependencies.CancelWithError(err)
			return
		}
	}

	// TODO move to env
	if _, err := debug.MakeContext(
		dependencies.Context,
		dependencies.Debug,
	); err != nil {
		dependencies.CancelWithError(err)
		return
	}

	env := env.Make(
		dependencies.Context,
		cmd.GetFlagSet(),
		dependencies.Config,
		dirLayout,
	)

	cmdArgs := cmd.Args()

	var u *repo_local.Repo

	options := repo_local.OptionsEmpty

	if og, ok := cmd.Command.(repo_local.OptionsGetter); ok {
		options = og.GetEnvironmentInitializeOptions()
	}

	{
		var err error

		if u, err = repo_local.Make(
			env,
			options,
		); err != nil {
			dependencies.CancelWithError(err)
			return
		}

		defer errors.DeferredFlusher(&err, u)
	}

	defer func() {
		if err := u.GetRepoLayout().ResetTempOnExit(
			dependencies.Context,
		); err != nil {
			dependencies.CancelWithError(err)
			return
		}
	}()

	switch {
	case u.GetConfig().Complete:
		var t CommandCompletionWithRepo
		haystack := any(cmd.Command)

	LOOP:
		for {
			switch c := haystack.(type) {
			case *commandWithQuery:
				t = c
				break LOOP

			case CommandCompletionWithRepo:
				t = c
				break LOOP

			default:
				dependencies.Cancel(errors.BadRequestf("Command does not support completion: %T", c))
				return
			}
		}

		t.CompleteWithRepo(u, cmdArgs...)

	default:

		func() {
			defer func() {
				// if r := recover(); r != nil {
				// 	result = ErrorResult{error: errors.Errorf("panicked: %s", r)}
				// }
			}()

			cmd.Command.RunWithRepo(u, cmdArgs...)
		}()
	}

	return
}
