package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
)

func (repo *Repo) MakeImporter(
	options sku.ImporterOptions,
	storeOptions sku.StoreOptions,
) (importer sku.Importer) {
	return repo.GetStore().MakeImporter(options, storeOptions)
}
