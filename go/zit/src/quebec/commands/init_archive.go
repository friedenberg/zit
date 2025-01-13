package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
)

type InitArchive struct {
	*flag.FlagSet
	repo_layout.BigBang
}

func init() {
	registerCommand(
		"init-archive",
		func(f *flag.FlagSet) Command {
			c := &InitArchive{
				FlagSet: f,
				BigBang: repo_layout.BigBang{},
			}

			c.SetFlagSet(f)
			c.Config.RepoType = repo_type.TypeArchive

			return c
		},
	)
}

func (c InitArchive) GetFlagSet() *flag.FlagSet {
	return c.FlagSet
}

func (c InitArchive) Run(
	dependencies Dependencies,
) {
	layout := dir_layout.MakeDefault(
		dependencies.Context,
		dependencies.Debug,
	)

	env := env.Make(
		dependencies.Context,
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
