package command_components

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
)

type Env struct{}

func (cmd *Env) MakeEnv(req command.Request) env.LocalEnv {
	return cmd.MakeEnvWithOptions(
		req,
		env.Options{},
	)
}

func (cmd *Env) MakeEnvWithOptions(
	req command.Request,
	options env.Options,
) env.LocalEnv {
	layout := env_dir.MakeDefault(
		req,
		req.Config.Debug,
	)

	return env.Make(
		req,
		req.Config,
		layout,
		options,
	)
}

func (cmd *Env) MakeEnvWithXDGLayoutAndOptions(
	req command.Request,
	xdgDotenvPath string,
	options env.Options,
) env.LocalEnv {
	dirLayout := env_dir.MakeFromXDGDotenvPath(
		req,
		req.Config,
		xdgDotenvPath,
	)

	env := env.Make(
		req,
		req.Config,
		dirLayout,
		options,
	)

	return env
}
