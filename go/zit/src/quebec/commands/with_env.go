package commands

import (
	"flag"

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

func (cmd commandWithEnv) Run(
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
		}
	}

	env := env.Make(
		dependencies.Context,
		dependencies.Config,
		dirLayout,
	)

	cmdArgs := cmd.Args()

	defer env.MustWithContext(env.GetDirLayout().ResetTempOnExit)

	cmd.Command.RunWithEnv(env, cmdArgs...)
}
