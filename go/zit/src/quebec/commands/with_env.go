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
	dirLayout := dir_layout.MakeDefault(
		dependencies.Context,
		dependencies.Debug,
	)

	var options env.Options

	if og, ok := cmd.Command.(env.OptionsGetter); ok {
		options = og.GetEnvOptions()
	}

	env := env.Make(
		dependencies.Context,
		dependencies.Config,
		dirLayout,
		options,
	)

	cmdArgs := cmd.Args()

	cmd.Command.RunWithEnv(env, cmdArgs...)
}
