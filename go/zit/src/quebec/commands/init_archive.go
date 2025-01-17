package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
)

func init() {
	command.Register(
		"init-archive",
		&InitArchive{
			BigBang: env_repo.BigBang{},
		},
	)
}

type InitArchive struct {
	env_repo.BigBang
}

func (c *InitArchive) SetFlagSet(f *flag.FlagSet) {
	c.BigBang.SetFlagSet(f)
	c.Config.RepoType = repo_type.TypeArchive
}

func (c InitArchive) Run(dependencies command.Request) {
	dir := env_dir.MakeDefault(
		dependencies,
		dependencies.Debug,
	)

	ui := env_ui.Make(
		dependencies,
		dependencies.Config,
		env_ui.Options{},
	)

	var repoLayout env_repo.Env

	layoutOptions := env_repo.Options{
		BasePath:             dependencies.Config.BasePath,
		PermitNoZitDirectory: true,
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

	repoLayout.Genesis(c.BigBang)
}
