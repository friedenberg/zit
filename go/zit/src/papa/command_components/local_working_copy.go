package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type LocalWorkingCopy struct {
	Env
}

func (cmd *LocalWorkingCopy) SetFlagSet(f *flag.FlagSet) {
}

func (c LocalWorkingCopy) MakeLocalWorkingCopy(
	dep command.Request,
) *local_working_copy.Repo {
	return c.MakeLocalWorkingCopyWithOptions(
		dep,
		env.Options{},
		local_working_copy.OptionsEmpty,
	)
}

func (cmd LocalWorkingCopy) MakeLocalWorkingCopyWithOptions(
	dep command.Request,
	envOptions env.Options,
	repoOptions local_working_copy.Options,
) *local_working_copy.Repo {
	env := cmd.MakeEnvWithOptions(dep, envOptions)

	return local_working_copy.Make(env, repoOptions)
}

// TODO modify to work with archives
func (cmd LocalWorkingCopy) MakeLocalWorkingCopyFromConfigAndXDGDotenvPath(
	req command.Request,
	xdgDotenvPath string,
	options env.Options,
) (local *local_working_copy.Repo) {
	env := cmd.MakeEnvWithXDGLayoutAndOptions(
		req,
		xdgDotenvPath,
		options,
	)

	local = local_working_copy.Make(
		env,
		local_working_copy.OptionsEmpty,
	)

	return
}
