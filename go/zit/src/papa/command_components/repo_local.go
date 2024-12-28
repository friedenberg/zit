package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type RepoLocal struct{}

func (cmd *RepoLocal) SetFlagSet(f *flag.FlagSet) {
}

func (c RepoLocal) MakeRepoLocal(
	context errors.Context,
	config config_mutable_cli.Config,
	options repo_local.Options,
) *repo_local.Repo {
	var layout dir_layout.Layout

	{
		var err error

		if layout, err = dir_layout.MakeDefault(
			config.Debug,
		); err != nil {
			context.CancelWithError(err)
		}
	}

	env := env.Make(
		context,
		config,
		layout,
	)

	return repo_local.Make(env, options)
}
