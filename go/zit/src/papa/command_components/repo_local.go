package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type RepoLocal struct{}

func (cmd *RepoLocal) SetFlagSet(f *flag.FlagSet) {
}

func (c RepoLocal) MakeRepoLocal(
	context *errors.Context,
	config config_mutable_cli.Config,
	envOptions env.Options,
	repoOptions local_working_copy.Options,
) *local_working_copy.Repo {
	layout := dir_layout.MakeDefault(
		context,
		config.Debug,
	)

	env := env.Make(
		context,
		config,
		layout,
		envOptions,
	)

	return local_working_copy.Make(env, repoOptions)
}
