package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/mike/importer"
)

func (repo *Repo) MakeImporter(
	options importer.ImporterOptions,
	storeOptions sku.StoreOptions,
) (importer importer.Importer) {
	return repo.GetStore().MakeImporter(options, storeOptions)
}
