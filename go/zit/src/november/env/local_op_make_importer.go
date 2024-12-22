package env

import (
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (env *Local) MakeImporter(
	printCopies bool,
) (importer store.Importer) {
	importer = store.Importer{
		Store: env.GetStore(),
	}

	if printCopies {
		importer.BlobCopierDelegate = func(result store.BlobCopyResult) error {
			// TODO switch to Err and fix test
			return ui.Out().Printf(
				"copied Blob %s (%d bytes)",
				result.GetBlobSha(),
				result.N,
			)
		}
	}

	return
}
