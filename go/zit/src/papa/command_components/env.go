package command_components

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
)

type Env struct{}

func (cmd *Env) MakeEnv(req command.Request) *env.Env {
	return cmd.MakeEnvWithOptions(
		req,
		env.Options{},
	)
}

func (cmd *Env) MakeEnvWithOptions(
	req command.Request,
	options env.Options,
) *env.Env {
	layout := dir_layout.MakeDefault(
		req.Context,
		req.Config.Debug,
	)

	return env.Make(
		req.Context,
		req.Config,
		layout,
		options,
	)
}

func (cmd *Env) MakeEnvWithXDGLayoutAndOptions(
	req command.Request,
	xdgDotenvPath string,
	options env.Options,
) *env.Env {
	dirLayout := dir_layout.MakeFromXDGDotenvPath(
		req.Context,
		req.Config,
		xdgDotenvPath,
	)

	env := env.Make(
		req.Context,
		req.Config,
		dirLayout,
		options,
	)

	return env
}
