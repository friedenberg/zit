package repo_local

import (
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (env *Repo) MakeImporter(
	printCopies bool,
) (importer store.Importer) {
	importer = store.Importer{
		Store: env.GetStore(),
	}

	if printCopies {
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
