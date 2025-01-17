package command_components

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
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
	layout := dir_layout.MakeDefault(
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
	dirLayout := dir_layout.MakeFromXDGDotenvPath(
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
