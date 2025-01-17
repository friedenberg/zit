package command_components

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
)

type Env struct{}

func (cmd *Env) MakeEnv(req command.Request) env_local.Env {
	return cmd.MakeEnvWithOptions(
		req,
		env_ui.Options{},
	)
}

func (cmd *Env) MakeEnvWithOptions(
	req command.Request,
	options env_ui.Options,
) env_local.Env {
	layout := env_dir.MakeDefault(
		req,
		req.Config.Debug,
	)

	return env_local.Make(
		env_ui.Make(
			req,
			req.Config,
			options,
		),
		layout,
	)
}

func (cmd *Env) MakeEnvWithXDGLayoutAndOptions(
	req command.Request,
	xdgDotenvPath string,
	options env_ui.Options,
) env_local.Env {
	dir := env_dir.MakeFromXDGDotenvPath(
		req,
		req.Config,
		xdgDotenvPath,
	)

	ui := env_ui.Make(
		req,
		req.Config,
		options,
	)

	return env_local.Make(ui, dir)
}
