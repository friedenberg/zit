package commands

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local_working_copy"
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

	repoOptions := repo_local_working_copy.OptionsEmpty

	if og, ok := cmd.Command.(repo_local_working_copy.OptionsGetter); ok {
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
