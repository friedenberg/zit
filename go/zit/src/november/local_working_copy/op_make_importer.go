package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (repo *Repo) MakeImporter(
	options store.ImporterOptions,
	storeOptions sku.StoreOptions,
) (importer store.Importer) {
	return repo.GetStore().MakeImporter(options, storeOptions)
}
