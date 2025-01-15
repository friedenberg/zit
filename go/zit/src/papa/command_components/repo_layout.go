package command_components

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
)

type RepoLayout struct{}

func (cmd RepoLayout) MakeRepoLayout(
	context *errors.Context,
	config config_mutable_cli.Config,
) repo_layout.Layout {
	layout := dir_layout.MakeDefault(
		context,
		config.Debug,
	)

	env := env.Make(
		context,
		config,
		layout,
		env.Options{},
	)

	var repoLayout repo_layout.Layout

	layoutOptions := repo_layout.Options{
		BasePath:             config.BasePath,
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

	return repoLayout
}
