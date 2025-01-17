package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
)

func init() {
	command.Register(
		"init-archive",
		&InitArchive{
			BigBang: repo_layout.BigBang{},
		},
	)
}

type InitArchive struct {
	repo_layout.BigBang
}

func (c *InitArchive) SetFlagSet(f *flag.FlagSet) {
	c.BigBang.SetFlagSet(f)
	c.Config.RepoType = repo_type.TypeArchive
}

func (c InitArchive) Run(dependencies command.Request) {
	layout := env_dir.MakeDefault(
		dependencies,
		dependencies.Debug,
	)

	env := env.Make(
		dependencies,
		dependencies.Config,
		layout,
		env.Options{},
	)

	var repoLayout repo_layout.Layout

	layoutOptions := repo_layout.Options{
		BasePath:             dependencies.Config.BasePath,
		PermitNoZitDirectory: true,
	}

	{
		var err error

		if repoLayout, err = repo_layout.Make(
			env,
			layoutOptions,
		); err != nil {
			env.CancelWithError(err)
		}

	}

	repoLayout.Genesis(c.BigBang)
}
