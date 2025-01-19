package store

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/mike/importer"
)

func (store *Store) MakeImporter(
	options importer.ImporterOptions,
	storeOptions sku.StoreOptions,
) importer.Importer {
	importer := importer.Make(
		options,
		storeOptions,
		store.GetEnvRepo(),
		store.GetTypedBlobStore().GetInventoryList(),
		store.GetStreamIndex(),
		store.externalStores[ids.RepoId{}],
		store,
	)

	return importer
}
