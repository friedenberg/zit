package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
)

type RepoLayout struct{}

func (cmd *RepoLayout) SetFlagSet(f *flag.FlagSet) {}

func (cmd RepoLayout) MakeRepoLayout(
	dep command.Request,
	permitNoZitDirectory bool,
) repo_layout.Layout {
	layout := env_dir.MakeDefault(
		dep,
		dep.Config.Debug,
	)

	env := env.Make(
		dep,
		dep.Config,
		layout,
		env.Options{},
	)

	var repoLayout repo_layout.Layout

	layoutOptions := repo_layout.Options{
		BasePath:             dep.Config.BasePath,
		PermitNoZitDirectory: permitNoZitDirectory,
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

	return repoLayout
}
