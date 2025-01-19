package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/repo_type"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_repo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/lima/inventory_list_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type LocalArchive struct{}

func (cmd *LocalArchive) SetFlagSet(f *flag.FlagSet) {
}

func (c LocalArchive) MakeLocalArchive(
	envRepo env_repo.Env,
) repo.Repo {
	repoType := envRepo.GetConfig().GetRepoType()

	switch repoType {
	case repo_type.TypeArchive:
		objectFormat := object_inventory_format.FormatForVersion(envRepo.GetStoreVersion())
		boxFormat := box_format.MakeBoxTransactedArchive(
			envRepo,
			options_print.V0{}.WithPrintTai(true),
		)

		inventoryListBlobStore := typed_blob_store.MakeInventoryStore(
			envRepo,
			objectFormat,
			boxFormat,
		)

		var inventoryListStore inventory_list_store.Store

		if err := inventoryListStore.Initialize(
			envRepo,
			objectFormat,
			nil,
			inventoryListBlobStore,
		); err != nil {
			envRepo.CancelWithError(err)
		}

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
