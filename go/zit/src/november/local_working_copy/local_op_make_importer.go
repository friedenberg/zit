package local_working_copy

import (
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (env *Repo) MakeImporter(
	options store.ImporterOptions,
) (importer store.Importer) {
	importer = store.Importer{
		Store:               env.GetStore(),
		ExcludeObjects:      options.ExcludeObjects,
		RemoteBlobStore:     options.RemoteBlobStore,
		BlobCopierDelegate:  options.BlobCopierDelegate,
		AllowMergeConflicts: options.AllowMergeConflicts,
		ParentNegotiator:    options.ParentNegotiator,
		CheckedOutPrinter:   options.CheckedOutPrinter,
	}

	if importer.BlobCopierDelegate == nil &&
		importer.RemoteBlobStore != nil &&
		options.PrintCopies {
		importer.BlobCopierDelegate = env.MakeBlobCopierDelegate()
	}

	return
}

func (env *Repo) MakeBlobCopierDelegate() func(store.BlobCopyResult) error {
	return func(result store.BlobCopyResult) error {
		return env.GetUI().Printf(
			"copied Blob %s (%d bytes)",
			result.Sha,
			result.N,
		)
	}
}
