package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/config_mutable_cli"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type BlobStoreLocal struct{}

func (cmd *BlobStoreLocal) SetFlagSet(f *flag.FlagSet) {
}

type BlobStoreWithEnv struct {
	env.Env
	interfaces.BlobStore
}

func (c BlobStoreLocal) MakeBlobStoreLocal(
	context errors.Context,
	config config_mutable_cli.Config,
	envOptions env.Options,
	repoOptions local_working_copy.Options,
) BlobStoreWithEnv {
	layout := env_dir.MakeDefault(
		context,
		config.Debug,
	)

	env := env.Make(
		context,
		config,
		layout,
		envOptions,
	)

	layoutOptions := repo_layout.Options{
		BasePath: config.BasePath,
	}

	var repoLayout repo_layout.Layout

	{
		var err error

		if repoLayout, err = repo_layout.Make(
			env,
			layoutOptions,
		); err != nil {
			context.CancelWithError(err)
		}
	}

	return BlobStoreWithEnv{
		Env:       env,
		BlobStore: repoLayout,
	}
}
