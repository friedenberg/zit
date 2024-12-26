package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
)

type commandWithEnv struct {
	Command CommandWithEnv
	*flag.FlagSet
}

func (cmd commandWithEnv) GetFlagSet() *flag.FlagSet {
	return cmd.FlagSet
}

func (cmd commandWithEnv) RunWithDependencies(
	dependencies Dependencies,
) (exitStatus int) {
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

	defer func() {
		if err := env.GetDirLayout().ResetTempOnExit(
			dependencies.Context,
		); err != nil {
			dependencies.CancelWithError(err)
			return
		}
	}()

	cmd.Command.RunWithEnv(env, cmdArgs...)

	return
}
