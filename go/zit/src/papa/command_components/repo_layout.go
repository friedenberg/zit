package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
)

type RepoLayout struct{}

func (cmd *RepoLayout) SetFlagSet(f *flag.FlagSet) {}

func (cmd RepoLayout) MakeRepoLayout(
	dep command.Request,
	permitNoZitDirectory bool,
) env_repo.Env {
	dir := env_dir.MakeDefault(
		dep,
		dep.Config.Debug,
	)

	ui := env_ui.Make(
		dep,
		dep.Config,
		env_ui.Options{},
	)

	var repoLayout env_repo.Env

	layoutOptions := env_repo.Options{
		BasePath:             dep.Config.BasePath,
		PermitNoZitDirectory: permitNoZitDirectory,
	}

	{
		var err error

		if repoLayout, err = env_repo.Make(
			env_local.Make(ui, dir),
			layoutOptions,
		); err != nil {
			ui.CancelWithError(err)
		}
	}

	return repoLayout
}
