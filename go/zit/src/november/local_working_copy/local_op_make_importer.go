package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (repo *Repo) MakeImporter(
	options store.ImporterOptions,
	storeOptions sku.StoreOptions,
) (importer store.Importer) {
	importer = store.Importer{
		Store:               repo.GetStore(),
		ExcludeObjects:      options.ExcludeObjects,
		RemoteBlobStore:     options.RemoteBlobStore,
		BlobCopierDelegate:  options.BlobCopierDelegate,
		AllowMergeConflicts: options.AllowMergeConflicts,
		ParentNegotiator:    options.ParentNegotiator,
		CheckedOutPrinter:   options.CheckedOutPrinter,
    StoreOptions: storeOptions,
	}

	if importer.BlobCopierDelegate == nil &&
		importer.RemoteBlobStore != nil &&
		options.PrintCopies {
		importer.BlobCopierDelegate = sku.MakeBlobCopierDelegate(repo.GetUI())
	}

	return
}
