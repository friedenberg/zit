package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/read_write_repo_local"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type commandWithBlobStore struct {
	*flag.FlagSet
	command_components.BlobStoreLocal
	Command CommandWithBlobStore
}

func (cmd commandWithBlobStore) GetFlagSet() *flag.FlagSet {
	return cmd.FlagSet
}

func (cmd commandWithBlobStore) Run(
	dependencies Dependencies,
) {
	var envOptions env.Options

	if og, ok := cmd.Command.(env.OptionsGetter); ok {
		envOptions = og.GetEnvOptions()
	}

	repoOptions := read_write_repo_local.OptionsEmpty

	if og, ok := cmd.Command.(read_write_repo_local.OptionsGetter); ok {
		repoOptions = og.GetLocalRepoOptions()
	}

	cmdArgs := cmd.Args()

	blobStore := cmd.MakeBlobStoreLocal(
		dependencies.Context,
		dependencies.Config,
		envOptions,
		repoOptions,
	)

	switch {
	case dependencies.Config.Complete:
		dependencies.CancelWithBadRequestf("completion not supported")

	default:
		cmd.Command.RunWithBlobStore(blobStore, cmdArgs...)
	}
}
