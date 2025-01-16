package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type LocalWorkingCopy struct{}

func (cmd *LocalWorkingCopy) SetFlagSet(f *flag.FlagSet) {
}

func (c LocalWorkingCopy) MakeLocalWorkingCopy(
	dep command.Dep,
	envOptions env.Options,
	repoOptions local_working_copy.Options,
) *local_working_copy.Repo {
	layout := dir_layout.MakeDefault(
		dep.Context,
		dep.Config.Debug,
	)

	env := env.Make(
		dep.Context,
		dep.Config,
		layout,
		envOptions,
	)

	return local_working_copy.Make(env, repoOptions)
}

// TODO modify to work with archives
func (cmd LocalWorkingCopy) MakeFromConfigAndXDGDotenvPath(
	context *errors.Context,
	config config_mutable_cli.Config,
	xdgDotenvPath string,
	options env.Options,
) (local *local_working_copy.Repo, err error) {
	dirLayout := dir_layout.MakeFromXDGDotenvPath(
		context,
		config,
		xdgDotenvPath,
	)

	env := env.Make(
		context,
		config,
		dirLayout,
		options,
	)

	local = local_working_copy.Make(
		env,
		local_working_copy.OptionsEmpty,
	)

	return
}
