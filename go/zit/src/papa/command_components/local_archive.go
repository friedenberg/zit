package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/lima/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/lima/inventory_list_store"
)

type LocalArchive struct{}

func (cmd *LocalArchive) SetFlagSet(f *flag.FlagSet) {
}

func (c LocalArchive) MakeLocalArchive(
	repoLayout repo_layout.Layout,
) *inventory_list_store.Store {
	objectFormat := object_inventory_format.FormatForVersion(repoLayout.GetStoreVersion())
	boxFormat := box_format.MakeBoxTransactedArchive(
		repoLayout.Env,
		options_print.V0{}.WithPrintTai(true),
	)

	inventoryListBlobStore := blob_store.MakeInventoryStore(
		repoLayout,
		objectFormat,
		boxFormat,
	)

	var inventoryListStore inventory_list_store.Store

	if err := inventoryListStore.Initialize(
		repoLayout,
		objectFormat,
		nil,
		inventoryListBlobStore,
	); err != nil {
		repoLayout.CancelWithError(err)
	}

	return &inventoryListStore
}
