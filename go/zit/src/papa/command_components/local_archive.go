package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/lima/inventory_list_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/mike/env_box"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type LocalArchive struct {
	EnvRepo
}

func (cmd *LocalArchive) SetFlagSet(f *flag.FlagSet) {
}

func (c LocalArchive) MakeLocalArchive(
	envRepo env_repo.Env,
) repo.LocalRepo {
	repoType := envRepo.GetConfig().ImmutableConfig.GetRepoType()

	switch repoType {
	case repo_type.TypeArchive:
		inventoryListBlobStore := c.MakeTypedInventoryListBlobStore(
			envRepo,
		)

		var inventoryListStore inventory_list_store.Store

		if err := inventoryListStore.Initialize(
			envRepo,
			nil,
			inventoryListBlobStore,
		); err != nil {
			envRepo.CancelWithError(err)
		}

		envBox := env_box.Make(
			envRepo,
			nil,
			nil,
		)

		inventoryListStore.SetUIDelegate(envBox.GetUIStorePrinters())

		return &inventoryListStore

	case repo_type.TypeWorkingCopy:
		return local_working_copy.MakeWithLayout(
			local_working_copy.OptionsEmpty,
			envRepo,
		)

	default:
		envRepo.CancelWithErrorf("unsupported repo type: %q", repoType)
		return nil
	}
}
